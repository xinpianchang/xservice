package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger log interface definition
type Logger interface {
	// Named see zap Named
	Named(name string) Logger

	// Debug log message
	Debug(msg string, fields ...zapcore.Field)

	// Info log message
	Info(msg string, fields ...zapcore.Field)

	// Warn log message
	Warn(msg string, fields ...zapcore.Field)

	// Error log message
	Error(msg string, fields ...zapcore.Field)

	// Fatal log message
	Fatal(msg string, fields ...zapcore.Field)

	// With add zap fields
	With(fields ...zapcore.Field) Logger

	// For log with context.Context, which will log trace_id and span_id if tracing enabled
	For(ctx context.Context) Logger

	// CallerSkip skip caller for adjust caller
	CallerSkip(int) Logger
}

// default logger implementation
type defaultLogger struct {
	logger           *zap.Logger
	additionalFields []zapcore.Field
	skiped           bool
}

// newLogger new logger instance
func newLogger(l *zap.Logger) Logger {
	return &defaultLogger{
		logger:           l,
		additionalFields: make([]zapcore.Field, 0, 8),
	}
}

// Named see zap Named
func (t defaultLogger) Named(name string) Logger {
	t.logger = t.logger.Named(name)
	return t
}

// Debug log message
func (t defaultLogger) Debug(msg string, fields ...zapcore.Field) {
	t.logger.Debug(msg, fields...)
}

// Info log message
func (t defaultLogger) Info(msg string, fields ...zapcore.Field) {
	t.logger.Info(msg, fields...)
}

// Warn log message
func (t defaultLogger) Warn(msg string, fields ...zapcore.Field) {
	t.logger.Warn(msg, fields...)
}

// Error log message
func (t defaultLogger) Error(msg string, fields ...zapcore.Field) {
	t.logger.Error(msg, fields...)
}

// Fatal log message
func (t defaultLogger) Fatal(msg string, fields ...zapcore.Field) {
	t.logger.Fatal(msg, fields...)
}

// With add zap fields
func (t defaultLogger) With(fields ...zapcore.Field) Logger {
	return defaultLogger{
		logger:           t.logger.WithOptions(zap.AddCallerSkip(t.skip())).With(fields...),
		additionalFields: append(t.additionalFields, fields...),
		skiped:           true,
	}
}

// For log with context.Context, which will log trace_id and span_id if tracing enabled
func (t defaultLogger) For(ctx context.Context) Logger {
	if ctx == nil {
		ctx = context.Background()
	}

	if span := trace.SpanFromContext(ctx); span != nil {
		l := spanLogger{span: span, logger: t.logger.WithOptions(zap.AddCallerSkip(t.skip())), additionalFields: t.additionalFields}

		l.logger = l.logger.With(
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)

		return l
	}

	return defaultLogger{logger: t.logger.WithOptions(zap.AddCallerSkip(-1))}
}

// CallerSkip skip caller for adjust caller
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
