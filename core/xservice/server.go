package xservice

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	gwrt "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	echomd "github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/core/middleware"
	"github.com/xinpianchang/xservice/pkg/echox"
	"github.com/xinpianchang/xservice/pkg/grpcx"
	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/signalx"
	"github.com/xinpianchang/xservice/pkg/tracingx"
)

// Server is the interface for xservice server
type Server interface {
	// Echo returns the echo instance
	Echo() *echo.Echo

	// Serve start and listen server
	Serve() error

	// GrpcRegister registers a gRPC service
	GrpcRegister(desc *grpc.ServiceDesc, impl interface{}, handler ...GrpcRegisterHandler)
}

type grpcService struct {
	Desc    *grpc.ServiceDesc
	Impl    interface{}
	Handler GrpcRegisterHandler
}

// GrpcRegisterHandler is the interface for gRPC service register handler
type GrpcRegisterHandler func(ctx context.Context, mux *gwrt.ServeMux, conn *grpc.ClientConn) error

type serverImpl struct {
	options      *Options
	grpcGateway  *gwrt.ServeMux
	echo         *echo.Echo
	grpc         *grpc.Server
	grpcServices []*grpcService
	httpHandler  http.Handler
}

func newServer(opts *Options) Server {
	server := &serverImpl{
		grpcServices: make([]*grpcService, 0, 128),
	}
	server.options = opts

	server.initEcho()
	server.initGrpc()

	return server
}

func (t *serverImpl) Echo() *echo.Echo {
	return t.echo
}

// Serve starts and listen server
func (t *serverImpl) Serve() error {
	address := t.getHttpAddress()

	log.Debug("serve", zap.String("address", address))

	upg, err := tableflip.New(tableflip.Options{
		UpgradeTimeout: time.Minute,
	})
	if err != nil {
		log.Fatal("tableflip init", zap.Error(err))
	}
	defer upg.Stop()

	t.waitSignalForTableflip(upg)

	ln, err := upg.Fds.Listen("tcp", address)
	if err != nil {
		log.Fatal("listen", zap.Error(err))
	}
	defer ln.Close()

	mux := cmux.New(ln)
	defer mux.Close()

	grpcL := mux.Match(cmux.HTTP2())
	defer grpcL.Close()

	httpL := mux.Match(cmux.HTTP1Fast())
	defer httpL.Close()

	if len(t.grpcServices) > 0 {
		go t.serveGrpc(grpcL)
	}

	server := http.Server{
		Handler:           t.httpHandler,
		ReadHeaderTimeout: time.Second * 30,
		IdleTimeout:       time.Minute * 1,
	}

	go func() {
		if err := server.Serve(httpL); err != nil {
			if err != http.ErrServerClosed && err != cmux.ErrServerClosed {
				log.Fatal("start http server", zap.Error(err))
			}
		}
	}()

	go func() {
		_ = mux.Serve()
	}()

	if err = upg.Ready(); err != nil {
		log.Fatal("ready", zap.Error(err))
	}

	// all ready
	t.registerGrpcServiceEtcd()

	signalx.AddShutdownHook(func(os.Signal) {
		_ = server.Shutdown(context.Background())
		t.grpc.GracefulStop()
		sentry.Flush(time.Second * 2)
		log.Info("shutdown", zap.Int("pid", os.Getpid()))
	})

	<-upg.Exit()

	signalx.Shutdown()

	return nil
}

// GrpcRegister registers a gRPC service
func (t *serverImpl) GrpcRegister(desc *grpc.ServiceDesc, impl interface{}, hs ...GrpcRegisterHandler) {
	var handler GrpcRegisterHandler
	if len(hs) > 0 {
		handler = hs[0]
	}
	t.grpcServices = append(t.grpcServices, &grpcService{desc, impl, handler})
}

func (t *serverImpl) getHttpAddress() string {
	return t.options.Config.GetString(core.ConfigServiceAddr)
}

func (t *serverImpl) waitSignalForTableflip(upg *tableflip.Upgrader) {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGUSR2, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		for s := range sig {
			switch s {
			case syscall.SIGUSR2, syscall.SIGHUP:
				err := upg.Upgrade()
				if err != nil {
					log.Error("upgrade failed", zap.Error(err))
					continue
				}
				log.Info("upgrade succeeded", zap.Int("pid", os.Getpid()))
				return
			default:
				upg.Stop()
			}
		}
	}()
}

func (t *serverImpl) initEcho() {
	e := t.newEcho("http")
	echox.ConfigValidator(e)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.Group("/debug/*", middleware.Pprof())

	t.echo = e
}

// init grpc
// add middleware https://github.com/grpc-ecosystem/go-grpc-middleware
func (t *serverImpl) initGrpc() {
	grpc.EnableTracing = true

	options := make([]grpc.ServerOption, 0, 8)
	options = append(options,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpcx.EnvoyproxyValidatorStreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcx.EnvoyproxyValidatorUnaryServerInterceptor(),
		)),
	)
	options = append(options, t.options.GrpcServerOptions...)
	g := grpc.NewServer(options...)
	t.grpc = g

	healthServer := health.NewServer()
	healthServer.SetServingStatus(t.options.Name, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(g, healthServer)

	t.grpcGateway = gwrt.NewServeMux(gwrt.WithRoutingErrorHandler(
		func(ctx context.Context, mux *gwrt.ServeMux, m gwrt.Marshaler, w http.ResponseWriter, r *http.Request, status int) {
			switch status {
			case http.StatusNotFound:
				t.echo.ServeHTTP(w, r)
			default:
				gwrt.DefaultRoutingErrorHandler(ctx, mux, m, w, r, status)
			}
		},
	))

	// echo instance for grpc-gateway, which wrap another echo instance, for gRPC service not found fallback serve
	e := t.newEcho("grpc_gateway")
	e.Use(echo.WrapMiddleware(func(handler http.Handler) http.Handler {
		return t.grpcGateway
	}))

	t.httpHandler = e
}

func (t *serverImpl) newEcho(subsystem string) *echo.Echo {
	e := echo.New()

	e.Logger = log.NewEchoLogger()
	e.IPExtractor = echo.ExtractIPFromXFFHeader(echo.TrustPrivateNet(true))
	e.HTTPErrorHandler = echox.HTTPErrorHandler

	// recover
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if x := recover(); x != nil {
					log.For(c.Request().Context()).Error("server panic error", zap.Any("error", x))
					_ = c.String(http.StatusInternalServerError, fmt.Sprint("internal server error, ", x))
				}
			}()
			return next(&echoContext{c})
		}
	})

	e.Use(echomd.RequestID())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(middleware.Trace(t.options.Config.GetBool("jaeger.body_dump"), t.options.EchoTracingSkipper))
	e.Use(middleware.Prometheus(strings.ReplaceAll(t.options.Name, "-", "_"), subsystem))

	// logger id & traceId & server-info
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Service", fmt.Sprint(t.options.Name, "/", t.options.Version, "/", t.options.Build))
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}
			c.Set(echo.HeaderXRequestID, id)
			ctx := context.WithValue(c.Request().Context(), core.ContextHeaderXRequestID, id)
			c.SetRequest(c.Request().WithContext(ctx))

			traceId := tracingx.GetTraceID(c.Request().Context())
			if traceId != "" {
				c.Response().Header().Set("X-Trace-Id", traceId)
			}

			if span := opentracing.SpanFromContext(c.Request().Context()); span != nil {
				span.SetTag("requestId", id)
				span.SetTag("ip", c.RealIP())
			}

			if hub := sentryecho.GetHubFromContext(c); hub != nil {
				scope := hub.Scope()
				scope.SetTag("ip", c.RealIP())
				scope.SetTag("X-Forwarded-For", c.Request().Header.Get("X-Forwarded-For"))
				if traceId != "" {
					scope.SetTag("traceId", traceId)
				}
			}

			return next(c)
		}
	})

	return e
}

func (t *serverImpl) serveGrpc(ln net.Listener) {
	for _, service := range t.grpcServices {
		t.grpc.RegisterService(service.Desc, service.Impl)
		// log.Debug("register grpc service", zap.String("impl", reflect.TypeOf(service.Impl).String()))
	}

	go func() {
		_ = t.grpc.Serve(ln)
	}()

	address := t.getHttpAddress()
	grpcClientConn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithInsecure(),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(),
			grpc_prometheus.StreamClientInterceptor,
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
			grpc_prometheus.UnaryClientInterceptor,
		)),
	)

	if err != nil {
		log.Fatal("grpc gateway client conn", zap.Error(err))
	}

	for _, service := range t.grpcServices {
		if service.Handler == nil {
			continue
		}
		err := service.Handler(context.Background(), t.grpcGateway, grpcClientConn)
		if err != nil {
			log.Fatal("grpc register handler", zap.Error(err))
		}
		// log.Debug("register grpc gateway", zap.String("handler", runtime.FuncForPC(reflect.ValueOf(service.Handler).Pointer()).Name()))
	}
}

// registerGrpcServiceEtcd
// refer: https://etcd.io/docs/v3.5/dev-guide/grpc_naming/
func (t *serverImpl) registerGrpcServiceEtcd() {
	if len(t.grpcServices) == 0 {
		return
	}

	if os.Getenv(core.EnvEtcd) == "" {
		log.Warn("etcd not configured, service register ignored")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	go t.doRegisterGrpcServiceEtcd(ctx)

	signalx.AddShutdownHook(func(s os.Signal) {
		cancel()
		// deregister
		log.Debug("deregister service")
		client := serviceEtcdClient()
		em, _ := endpoints.NewManager(client, core.ServiceRegisterKeyPrefix)

		for _, service := range t.grpcServices {
			_ = em.DeleteEndpoint(context.Background(), serviceKey(os.Getenv(core.EnvServiceName), service.Desc))
		}
	})
}

func (t *serverImpl) doRegisterGrpcServiceEtcd(ctx context.Context) {
	l := log.Named("registerGrpcServiceEtcd")
	defer func() {
		if x := recover(); x != nil {
			l.Error("recover", zap.Any("err", x))
			sentry.CaptureException(errors.WithStack(errors.New(fmt.Sprint(x))))

			time.Sleep(time.Second * 10)
			go t.doRegisterGrpcServiceEtcd(ctx)
		}
	}()

	client := serviceEtcdClient()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	ttl := int64(10) // seconds

	em, _ := endpoints.NewManager(client, core.ServiceRegisterKeyPrefix)

	var (
		id    clientv3.LeaseID = 0
		lease clientv3.Lease   = clientv3.NewLease(client)

		addr = t.options.Config.GetString(core.ConfigServiceAdvertisedAddr)
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if id == 0 {
				leaseRsp, err := lease.Grant(context.Background(), ttl)
				if err != nil {
					l.Error("lease.Grant", zap.Error(err))
					continue
				}
				id = leaseRsp.ID

				for _, service := range t.grpcServices {
					key := serviceKey(os.Getenv(core.EnvServiceName), service.Desc)
					endpoint := endpoints.Endpoint{
						Addr:     addr,
						Metadata: service.Desc.Metadata,
					}

					ll := l.With(zap.String("service", key))
					err = em.AddEndpoint(context.Background(), key, endpoint, clientv3.WithLease(id))
					if err != nil {
						ll.Error("kv.Put", zap.Error(err))
					}
				}
			} else {
				_, err := lease.KeepAliveOnce(context.Background(), id)
				if err != nil {
					id = 0
				}
			}
		}

		// wait next loop
		<-ticker.C
	}
}

type echoContext struct {
	echo.Context
}

func (t *echoContext) Logger() echo.Logger {
	logger := t.Context.Logger()
	if l, ok := logger.(*log.EchoLogger); ok {
		return l.For(t.Request().Context())
	}
	return logger
}

func (t *echoContext) Path() string {
	return t.Request().URL.Path
}
