package log

import (
	"context"
	"fmt"
	"math"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger           *zap.Logger
	span             trace.Span
	additionalFields []zapcore.Field
}

// Named name log
// see zap log named
func (t spanLogger) Named(name string) Logger {
	t.logger = t.logger.Named(name)
	return t
}

// Debug log message
func (t spanLogger) Debug(msg string, fields ...zapcore.Field) {
	t.logger.Debug(msg, fields...)
}

// Info log message
func (t spanLogger) Info(msg string, fields ...zapcore.Field) {
	t.logger.Info(msg, fields...)
}

// Warn log message
func (t spanLogger) Warn(msg string, fields ...zapcore.Field) {
	t.logToSpan("warn", msg, append(t.additionalFields, fields...)...)
	t.logger.Warn(msg, fields...)
}

// Error log message
func (t spanLogger) Error(msg string, fields ...zapcore.Field) {
	t.logToSpan("error", msg, fields...)
	t.logger.Error(msg, fields...)
}

// Fatal log message
func (t spanLogger) Fatal(msg string, fields ...zapcore.Field) {
	t.logToSpan("fatal", msg, append(t.additionalFields, fields...)...)
	t.logger.Fatal(msg, fields...)
}

// With add zap fields
func (t spanLogger) With(fields ...zapcore.Field) Logger {
	return spanLogger{
		logger:           t.logger.With(fields...),
		span:             t.span,
		additionalFields: append(t.additionalFields, fields...),
	}
}

// For log with context.Context, which will log trace_id and span_id if tracing enabled
func (t spanLogger) For(context.Context) Logger {
	return t
}

// CallerSkip skip caller for adjust caller
func (t spanLogger) CallerSkip(skip int) Logger {
	t.logger = t.logger.WithOptions(zap.AddCallerSkip(skip))
	return t
}

func (t spanLogger) logToSpan(level string, msg string, fields ...zapcore.Field) {
	fs := make([]attribute.KeyValue, 0, 2+len(fields))
	fs = append(fs, attribute.String("msg", msg))
	fs = append(fs, attribute.String("level", level))
	for _, field := range fields {
		fs = append(fs, zapFieldToKv(field))
	}
	t.span.SetAttributes(fs...)
}

// zapFieldToKv to tracing attribute
func zapFieldToKv(field zapcore.Field) attribute.KeyValue {
	switch field.Type {
	case zapcore.BoolType:
		val := false
		if field.Integer >= 1 {
			val = true
		}
		return attribute.Bool(field.Key, val)
	case zapcore.Float32Type:
		return attribute.Float64(field.Key, float64(math.Float32frombits(uint32(field.Integer))))
	case zapcore.Float64Type:
		return attribute.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.Int64Type:
		return attribute.Int64(field.Key, int64(field.Integer))
	case zapcore.Int32Type:
		return attribute.Int64(field.Key, int64(field.Integer))
	case zapcore.StringType:
		return attribute.String(field.Key, field.String)
	case zapcore.StringerType:
		return attribute.String(field.Key, field.Interface.(fmt.Stringer).String())
	case zapcore.Uint64Type:
		return attribute.String(field.Key, fmt.Sprint(field.Integer))
	case zapcore.Uint32Type:
		return attribute.String(field.Key, fmt.Sprint(field.Integer))
	case zapcore.DurationType:
		return attribute.String(field.Key, time.Duration(field.Integer).String())
	case zapcore.ErrorType:
		return attribute.String("err", (field.Interface.(error)).Error())
	default:
		return attribute.String(field.Key, fmt.Sprint(field.Interface))
	}
}
