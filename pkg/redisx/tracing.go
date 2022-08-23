package redisx

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/opentracing/opentracing-go"
)

type redisTracing struct{}

func (redisTracing) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if span := opentracing.SpanFromContext(ctx); span == nil {
		return ctx, nil
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, cmd.FullName())
	span.SetTag("db.system", "redis")
	span.SetTag("db.statement", cmd.String())

	return ctx, nil
}

func (redisTracing) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.Finish()
	}
	return nil
}

func (redisTracing) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if span := opentracing.SpanFromContext(ctx); span == nil {
		return ctx, nil
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "pipline")
	span.SetTag("db.system", "redis")
	span.SetTag("db.cmd_count", len(cmds))

	return ctx, nil
}

func (t redisTracing) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.Finish()
	}
	return nil
}
