package log

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger           *zap.Logger
	span             opentracing.Span
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

// For log with context.Context, which will log trace_id and span_id if opentracing enabled
func (t spanLogger) For(context.Context) Logger {
	return t
}

// CallerSkip skip caller for adjust caller
func (t spanLogger) CallerSkip(skip int) Logger {
	t.logger = t.logger.WithOptions(zap.AddCallerSkip(skip))
	return t
}

func (t spanLogger) logToSpan(level string, msg string, fields ...zapcore.Field) {
	fs := make([]log.Field, 0, 2+len(fields))
	fs = append(fs, log.String("msg", msg))
	fs = append(fs, log.String("level", level))
	for _, field := range fields {
		fs = append(fs, zapFieldToLogField(field))
	}
	t.span.LogFields(fs...)
}

// zapFieldToLogField to opentracing log field
func zapFieldToLogField(field zapcore.Field) log.Field {
	switch field.Type {
	case zapcore.BoolType:
		val := false
		if field.Integer >= 1 {
			val = true
		}
		return log.Bool(field.Key, val)
	case zapcore.Float32Type:
		return log.Float32(field.Key, math.Float32frombits(uint32(field.Integer)))
	case zapcore.Float64Type:
		return log.Float64(field.Key, math.Float64frombits(uint64(field.Integer)))
	case zapcore.Int64Type:
		return log.Int64(field.Key, int64(field.Integer))
	case zapcore.Int32Type:
		return log.Int32(field.Key, int32(field.Integer))
	case zapcore.StringType:
		return log.String(field.Key, field.String)
	case zapcore.StringerType:
		return log.String(field.Key, field.Interface.(fmt.Stringer).String())
	case zapcore.Uint64Type:
		return log.Uint64(field.Key, uint64(field.Integer))
	case zapcore.Uint32Type:
		return log.Uint32(field.Key, uint32(field.Integer))
	case zapcore.DurationType:
		return log.String(field.Key, time.Duration(field.Integer).String())
	case zapcore.ErrorType:
		return log.Error(field.Interface.(error))
	default:
		return log.Object(field.Key, field.Interface)
	}
}
