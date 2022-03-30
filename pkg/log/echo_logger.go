package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// EchoLogger echo logger
// see echo.Logger
type EchoLogger struct {
	prefix string
	logger Logger
}

// NewEchoLogger crate new echo log instance
func NewEchoLogger() *EchoLogger {
	return &EchoLogger{
		logger: Named(""),
	}
}

var _ echo.Logger = (*EchoLogger)(nil)

// For with context.Context
func (t *EchoLogger) For(ctx context.Context) *EchoLogger {
	return &EchoLogger{prefix: t.prefix, logger: t.logger.For(ctx)}
}

// Output get output
func (t *EchoLogger) Output() io.Writer {
	return nil
}

// SetOutput not implemented, ignored
func (t *EchoLogger) SetOutput(w io.Writer) {}

// Prefix get prefix
func (t *EchoLogger) Prefix() string {
	return t.prefix
}

// SetOutput set output prefix
func (t *EchoLogger) SetPrefix(p string) {
	t.prefix = p
	t.logger = Named(p)
}

// Level get log level
func (t *EchoLogger) Level() log.Lvl {
	return log.INFO
}

// SetLvl set log level
func (t *EchoLogger) SetLevel(v log.Lvl) {}

// SetHeader set header, not implemented
func (t *EchoLogger) SetHeader(h string) {}

// Print log
func (t *EchoLogger) Print(i ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Printf log formatted
func (t *EchoLogger) Printf(format string, args ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Printj log json formatted
func (t *EchoLogger) Printj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Info(strings.TrimSpace(string(b)))
}

// Debug log debug message
func (t *EchoLogger) Debug(i ...interface{}) {
	t.logger.Debug(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Debugf log debug formatted message
func (t *EchoLogger) Debugf(format string, args ...interface{}) {
	t.logger.Debug(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Debugj log json formatted
func (t *EchoLogger) Debugj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Debug(strings.TrimSpace(string(b)))
}

// Info log
func (t *EchoLogger) Info(i ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Infof log formatted
func (t *EchoLogger) Infof(format string, args ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Infoj log json
func (t *EchoLogger) Infoj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Info(strings.TrimSpace(string(b)))
}

// Warn log warn message
func (t *EchoLogger) Warn(i ...interface{}) {
	t.logger.Warn(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Warnf log warn formatted message
func (t *EchoLogger) Warnf(format string, args ...interface{}) {
	t.logger.Warn(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Warnj log warn json
func (t *EchoLogger) Warnj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Warn(strings.TrimSpace(string(b)))
}

// Error log error message
func (t *EchoLogger) Error(i ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Errorf log error formatted message
func (t *EchoLogger) Errorf(format string, args ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Errorj log error json
func (t *EchoLogger) Errorj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Error(strings.TrimSpace(string(b)))
}

// Fatal log fatal message
func (t *EchoLogger) Fatal(i ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Fatalj log fatal json message
func (t *EchoLogger) Fatalj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Error(strings.TrimSpace(string(b)))
}

// Fatalf log formatted message
func (t *EchoLogger) Fatalf(format string, args ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Panic log panic message, which will cause panic
func (t *EchoLogger) Panic(i ...interface{}) {
	t.logger.Fatal(strings.TrimSpace(fmt.Sprintln(i...)))
}

// Panicj log panic json message, which will cause panic
func (t *EchoLogger) Panicj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Fatal(strings.TrimSpace(string(b)))
}

// Panicf log formatted message, which will cause panic
func (t *EchoLogger) Panicf(format string, args ...interface{}) {
	t.logger.Fatal(strings.TrimSpace(fmt.Sprintf(format, args...)))
}
