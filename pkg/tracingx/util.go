package tracingx

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func GetTraceID(ctx context.Context) string {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			return sc.TraceID().String()
		}
	}

	return ""
}
