package tracingx

import (
	"context"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xinpianchang/xservice/v2/core"
	"github.com/xinpianchang/xservice/v2/pkg/log"
	"github.com/xinpianchang/xservice/v2/pkg/signalx"
)

// Config config tracing
func Config(v *viper.Viper) {
	serviceName := os.Getenv(core.EnvServiceName)

	if !v.IsSet("tracing") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),
	))

	if err != nil {
		log.Fatal("create tracing resource", zap.Error(err))
	}

	conn, err := grpc.DialContext(ctx, v.GetString("tracing.endpoint"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatal("init tracing", zap.Error(err))
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatal("create tracing exporter", zap.Error(err))
	}

	bsp := trace.NewBatchSpanProcessor(exporter)
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	signalx.AddShutdownHook(func(s os.Signal) {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Error("shutdown tracing", zap.Error(err))
		}
	})
}
