package log

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger           *zap.Logger
	span             opentracing.Span
	spanFields       []zapcore.Field
	additionalFields []zapcore.Field
}

func (t spanLogger) Named(name string) Logger {
	t.logger = t.logger.Named(name)
	return t
}

func (t spanLogger) Debug(msg string, fields ...zapcore.Field) {
	if t.logger.Core().Enabled(zap.DebugLevel) {
		t.logToSpan("debug", msg, append(t.additionalFields, fields...)...)
	}
	t.logger.Debug(msg, append(t.spanFields, fields...)...)
}

func (t spanLogger) Info(msg string, fields ...zapcore.Field) {
	if t.logger.Core().Enabled(zap.InfoLevel) {
		t.logToSpan("info", msg, append(t.additionalFields, fields...)...)
	}
	t.logger.Info(msg, append(t.spanFields, fields...)...)
}

func (t spanLogger) Warn(msg string, fields ...zapcore.Field) {
	t.logToSpan("warn", msg, append(t.additionalFields, fields...)...)
	t.logger.Warn(msg, append(t.spanFields, fields...)...)
}

func (t spanLogger) Error(msg string, fields ...zapcore.Field) {
	t.logToSpan("error", msg, fields...)
	ext.Error.Set(t.span, true)
	t.logger.Error(msg, append(t.spanFields, fields...)...)
}

func (t spanLogger) Fatal(msg string, fields ...zapcore.Field) {
	t.logToSpan("fatal", msg, append(t.additionalFields, fields...)...)
	ext.Error.Set(t.span, true)
	t.logger.Fatal(msg, append(t.spanFields, fields...)...)
}

func (t spanLogger) With(fields ...zapcore.Field) Logger {
	return spanLogger{
		logger:           t.logger.With(fields...),
		span:             t.span,
		spanFields:       t.spanFields,
		additionalFields: append(t.additionalFields, fields...),
	}
}

func (t spanLogger) For(context.Context) Logger {
	return t
}

func (t spanLogger) CallerSkip(skip int) Logger {
	t.logger = t.logger.WithOptions(zap.AddCallerSkip(skip))
	return t
}

func (t spanLogger) logToSpan(level string, msg string, fields ...zapcore.Field) {
	fa := fieldAdapter(make([]log.Field, 0, 2+len(fields)))
	fa = append(fa, log.String("msg", msg))
	fa = append(fa, log.String("level", level))
	for _, field := range fields {
		field.AddTo(&fa)
	}
	t.span.LogFields(fa...)
}

type fieldAdapter []log.Field

func (fa *fieldAdapter) AddBool(key string, value bool) {
	*fa = append(*fa, log.Bool(key, value))
}

func (fa *fieldAdapter) AddFloat64(key string, value float64) {
	*fa = append(*fa, log.Float64(key, value))
}

func (fa *fieldAdapter) AddFloat32(key string, value float32) {
	*fa = append(*fa, log.Float64(key, float64(value)))
}

func (fa *fieldAdapter) AddInt(key string, value int) {
	*fa = append(*fa, log.Int(key, value))
}

func (fa *fieldAdapter) AddInt64(key string, value int64) {
	*fa = append(*fa, log.Int64(key, value))
}

func (fa *fieldAdapter) AddInt32(key string, value int32) {
	*fa = append(*fa, log.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt16(key string, value int16) {
	*fa = append(*fa, log.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt8(key string, value int8) {
	*fa = append(*fa, log.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint(key string, value uint) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint64(key string, value uint64) {
	*fa = append(*fa, log.Uint64(key, value))
}

func (fa *fieldAdapter) AddUint32(key string, value uint32) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint16(key string, value uint16) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint8(key string, value uint8) {
	*fa = append(*fa, log.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUintptr(key string, value uintptr)                        {}
func (fa *fieldAdapter) AddArray(key string, marshaler zapcore.ArrayMarshaler) error { return nil }
func (fa *fieldAdapter) AddComplex128(key string, value complex128)                  {}
func (fa *fieldAdapter) AddComplex64(key string, value complex64)                    {}
func (fa *fieldAdapter) AddObject(key string, value zapcore.ObjectMarshaler) error   { return nil }
func (fa *fieldAdapter) AddReflected(key string, value interface{}) error            { return nil }
func (fa *fieldAdapter) OpenNamespace(key string)                                    {}

func (fa *fieldAdapter) AddDuration(key string, value time.Duration) {
	*fa = append(*fa, log.String(key, value.String()))
}

func (fa *fieldAdapter) AddTime(key string, value time.Time) {
	*fa = append(*fa, log.String(key, value.String()))
}

func (fa *fieldAdapter) AddBinary(key string, value []byte) {
	*fa = append(*fa, log.Object(key, value))
}

func (fa *fieldAdapter) AddByteString(key string, value []byte) {
	*fa = append(*fa, log.Object(key, value))
}

func (fa *fieldAdapter) AddString(key, value string) {
	if key != "" && value != "" {
		*fa = append(*fa, log.String(key, value))
	}
}
