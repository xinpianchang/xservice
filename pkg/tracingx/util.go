package tracingx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func GetTraceID(ctx context.Context) string {
	if span := trace.SpanFromContext(ctx); span != nil {
		if span.SpanContext().HasTraceID() {
			return span.SpanContext().TraceID().String()
		}
	}
	return ""
}
