package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/labstack/gommon/log"
)

type EchoLogger struct {
	prefix string
	logger Logger
}

func NewEchoLogger() *EchoLogger {
	return &EchoLogger{
		logger: Named(""),
	}
}

func (t *EchoLogger) For(ctx context.Context) *EchoLogger {
	return &EchoLogger{prefix: t.prefix, logger: t.logger.For(ctx)}
}

func (t *EchoLogger) Output() io.Writer {
	return nil
}

func (t *EchoLogger) SetOutput(w io.Writer) {}

func (t *EchoLogger) Prefix() string {
	return t.prefix
}

func (t *EchoLogger) SetPrefix(p string) {
	t.prefix = p
	t.logger = Named(p)
}

func (t *EchoLogger) Level() log.Lvl {
	return log.INFO
}

func (t *EchoLogger) SetLevel(v log.Lvl) {}

func (t *EchoLogger) SetHeader(h string) {}

func (t *EchoLogger) Print(i ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Printf(format string, args ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func (t *EchoLogger) Printj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Info(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Debug(i ...interface{}) {
	t.logger.Debug(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Debugf(format string, args ...interface{}) {
	t.logger.Debug(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func (t *EchoLogger) Debugj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Debug(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Info(i ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Infof(format string, args ...interface{}) {
	t.logger.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func (t *EchoLogger) Infoj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Info(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Warn(i ...interface{}) {
	t.logger.Warn(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Warnf(format string, args ...interface{}) {
	t.logger.Warn(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func (t *EchoLogger) Warnj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Warn(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Error(i ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Errorf(format string, args ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func (t *EchoLogger) Errorj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Error(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Fatal(i ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Fatalj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Error(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Fatalf(format string, args ...interface{}) {
	t.logger.Error(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func (t *EchoLogger) Panic(i ...interface{}) {
	t.logger.Fatal(strings.TrimSpace(fmt.Sprintln(i...)))
}

func (t *EchoLogger) Panicj(j log.JSON) {
	b, _ := json.Marshal(j)
	t.logger.Fatal(strings.TrimSpace(string(b)))
}

func (t *EchoLogger) Panicf(format string, args ...interface{}) {
	t.logger.Fatal(strings.TrimSpace(fmt.Sprintf(format, args...)))
}
