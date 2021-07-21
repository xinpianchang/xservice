package tracingx

import (
	"fmt"
	"os"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/signalx"
)

func Config(v *viper.Viper) {
	serviceName := os.Getenv(core.EnvServiceName)

	for k, v := range v.GetStringMapString("jaeger") {
		os.Setenv(fmt.Sprintf("JAEGER_%s", strings.ToUpper(k)), v)
	}

	os.Setenv("JAEGER_SERVICE_NAME", serviceName)

	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	if cfg.Sampler.Type == "" {
		cfg.Sampler.Type = "const"
	}
	if cfg.Sampler.Param == 0 {
		cfg.Sampler.Param = 1
	}

	metricsFactory := prometheus.New().Namespace(metrics.NSOptions{Name: "xservice", Tags: map[string]string{"service": serviceName}})

	tracer, closer, err := cfg.NewTracer(
		config.Logger(&jaegerLoggerAdapter{}),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)

	if err != nil {
		log.Fatal("create tracer", zap.Error(err))
	}

	signalx.AddShutdownHook(func(os.Signal) {
		_ = closer.Close()
	})

	opentracing.SetGlobalTracer(tracer)
}

type jaegerLoggerAdapter struct{}

func (l jaegerLoggerAdapter) Error(msg string) {
	log.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	log.Info(fmt.Sprintf(msg, args...))
}

func (l jaegerLoggerAdapter) Debugf(msg string, args ...interface{}) {
	// ignore debug
	// log.Debug(fmt.Sprintf(msg, args...))
}
