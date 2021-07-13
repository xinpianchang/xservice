package log

import (
	"context"

	"go.uber.org/zap"
)

// Named name log
// see zap log named
func Named(name string) Logger {
	return logger.Named(name)
}

// Debug log message
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info log message
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn log message
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error log message
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal log message
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// With log message
func With(fields ...zap.Field) Logger {
	return logger.With(fields...)
}

// For log with context.Context, which will log trace_id and span_id if opentracing enabled
func For(ctx context.Context) Logger {
	return logger.For(ctx)
}

// Get get zap log instance (global instance only)
func Get() *zap.Logger {
	return zaplogger
}
