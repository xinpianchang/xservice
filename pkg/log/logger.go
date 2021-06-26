package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Named(name string) Logger
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	With(fields ...zapcore.Field) Logger
	For(ctx context.Context) Logger
	CallerSkip(int) Logger
}

type defaultLogger struct {
	logger           *zap.Logger
	additionalFields []zapcore.Field
	skiped           bool
}

func newLogger(l *zap.Logger) Logger {
	return &defaultLogger{
		logger:           l,
		additionalFields: make([]zapcore.Field, 0, 8),
	}
}

func (t defaultLogger) Named(name string) Logger {
	t.logger = t.logger.Named(name)
	return t
}

func (t defaultLogger) Debug(msg string, fields ...zapcore.Field) {
	t.logger.Debug(msg, fields...)
}

func (t defaultLogger) Info(msg string, fields ...zapcore.Field) {
	t.logger.Info(msg, fields...)
}

func (t defaultLogger) Warn(msg string, fields ...zapcore.Field) {
	t.logger.Warn(msg, fields...)
}

func (t defaultLogger) Error(msg string, fields ...zapcore.Field) {
	t.logger.Error(msg, fields...)
}

func (t defaultLogger) Fatal(msg string, fields ...zapcore.Field) {
	t.logger.Fatal(msg, fields...)
}

func (t defaultLogger) With(fields ...zapcore.Field) Logger {
	return defaultLogger{
		logger:           t.logger.WithOptions(zap.AddCallerSkip(t.skip())).With(fields...),
		additionalFields: append(t.additionalFields, fields...),
		skiped:           true,
	}
}

func (t defaultLogger) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		l := spanLogger{span: span, logger: t.logger.WithOptions(zap.AddCallerSkip(t.skip())), additionalFields: t.additionalFields}

		if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
			l.spanFields = []zapcore.Field{
				zap.String("trace_id", jaegerCtx.TraceID().String()),
				zap.String("span_id", jaegerCtx.SpanID().String()),
			}
		}

		return l
	}

	return defaultLogger{logger: t.logger.WithOptions(zap.AddCallerSkip(-1))}
}

func (t defaultLogger) CallerSkip(skip int) Logger {
	t.logger = t.logger.WithOptions(zap.AddCallerSkip(skip))
	return t
}

func (t defaultLogger) skip() int {
	if t.skiped {
		return 0
	}
	return -1
}
