package log

import (
	"context"

	"go.uber.org/zap"
)

func Named(name string) Logger {
	return logger.Named(name)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func With(fields ...zap.Field) Logger {
	return logger.With(fields...)
}

func For(ctx context.Context) Logger {
	return logger.For(ctx)
}

func Get() *zap.Logger {
	return zaplogger
}
