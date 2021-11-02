package xservice

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	gresolver "google.golang.org/grpc/resolver"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/signalx"
)

// Client is the client for xservice
type Client interface {
	// GrpcClientConn returns a grpc client connection
	GrpcClientConn(ctx context.Context, service string, desc *grpc.ServiceDesc, endpoint ...string) (grpc.ClientConnInterface, error)
}

type clientImpl struct {
	options   *Options
	resolver  gresolver.Builder
	conn      map[string]*grpc.ClientConn
	connMutex sync.RWMutex
}

func newClient(opts *Options) Client {
	client := &clientImpl{
		options: opts,
		conn:    make(map[string]*grpc.ClientConn, 128),
	}

	if os.Getenv(core.EnvEtcd) != "" {
		var err error
		client.resolver, err = resolver.NewBuilder(serviceEtcdClient())
		if err != nil {
			log.Fatal("endpoints manager", zap.Error(err))
		}
	}

	signalx.AddShutdownHook(func(os.Signal) {
		client.connMutex.Lock()
		defer client.connMutex.Unlock()
		for _, c := range client.conn {
			_ = c.Close()
		}
	})

	return client
}

// GrpcClientConn returns a grpc client connection
func (t *clientImpl) GrpcClientConn(ctx context.Context, service string, desc *grpc.ServiceDesc, endpoint ...string) (grpc.ClientConnInterface, error) {
	key := t.grpcClientKey(service, desc, endpoint...)

	client := t.fastGetGrpcClient(key)
	if client != nil {
		return client, nil
	}

	t.connMutex.Lock()
	defer t.connMutex.Unlock()

	// double check
	if c, ok := t.conn[key]; ok {
		return c, nil
	}

	options := make([]grpc.DialOption, 0, 8)
	options = append(options,
		grpc.WithInsecure(), grpc.WithBlock(),
		// refer: https://github.com/grpc/grpc-go/blob/master/examples/features/load_balancing/client/main.go#L76
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)
	options = append(options,
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpc_opentracing.StreamClientInterceptor(),
			grpc_prometheus.StreamClientInterceptor,
		)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
			grpc_prometheus.UnaryClientInterceptor,
		)),
	)
	options = append(options, t.options.GrpcClientDialOptions...)

	ctx, cancel := context.WithTimeout(ctx, t.options.GrpcClientDialTimeout)
	defer cancel()

	if len(endpoint) > 0 {
		c, err := grpc.DialContext(ctx, endpoint[0], options...)
		if err != nil {
			return nil, err
		}

		t.conn[key] = c

		return c, nil
	}

	if os.Getenv(core.EnvEtcd) == "" {
		log.Fatal("etcd not configured")
	}

	target := fmt.Sprint("etcd:///", serviceKeyPrefix(service, desc))
	options = append(options, grpc.WithResolvers(t.resolver))
	c, err := grpc.DialContext(ctx, target, options...)
	if err != nil {
		return nil, err
	}

	t.conn[key] = c

	return c, nil
}

func (t *clientImpl) fastGetGrpcClient(key string) grpc.ClientConnInterface {
	t.connMutex.RLock()
	defer t.connMutex.RUnlock()
	if c, ok := t.conn[key]; ok {
		return c
	}
	return nil
}

func (t *clientImpl) grpcClientKey(service string, desc *grpc.ServiceDesc, endpoint ...string) string {
	var sb strings.Builder
	_, _ = sb.WriteString(service)
	_, _ = sb.WriteString(desc.ServiceName)
	if endpoint != nil {
		_, _ = sb.WriteString(endpoint[0])
	}
	return sb.String()
}
