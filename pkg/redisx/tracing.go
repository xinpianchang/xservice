package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type redisTracing struct{}

var tracer = otel.Tracer("redisx")

// DialHook implements redis.Hook
func (t *redisTracing) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook implements redis.Hook
func (*redisTracing) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if span := trace.SpanFromContext(ctx); span == nil {
			ctx, span = tracer.Start(ctx, cmd.FullName())
			span.SetAttributes(attribute.String("db.system", "redis"))
			span.SetAttributes(attribute.String("db.statement", cmd.String()))
			defer span.End()
		}
		return next(ctx, cmd)
	}
}

// ProcessPipelineHook implements redis.Hook
func (*redisTracing) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if span := trace.SpanFromContext(ctx); span == nil {
			ctx, span = tracer.Start(ctx, "pipline")
			span.SetAttributes(attribute.String("db.system", "redis"))
			span.SetAttributes(attribute.Int("db.cmd_count", len(cmds)))
			defer span.End()
		}
		return next(ctx, cmds)
	}
}

var _ redis.Hook = (*redisTracing)(nil)
