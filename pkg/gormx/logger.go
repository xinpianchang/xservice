package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xinpianchang/xservice/pkg/log"
)

type dbLogger struct {
	logger log.Logger
}

func (t *dbLogger) LogMode(logger.LogLevel) logger.Interface {
	return t
}

func (t *dbLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	t.logger.For(ctx).CallerSkip(2).Info(fmt.Sprintf(msg, data...))
}

func (t *dbLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	t.logger.For(ctx).CallerSkip(2).Warn(fmt.Sprintf(msg, data...))
}

func (t *dbLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	t.logger.For(ctx).CallerSkip(2).Error(fmt.Sprintf(msg, data...))
}

func (t *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	l := t.logger.For(ctx).CallerSkip(2)
	sql, rows := fc()
	l = l.With(zap.Duration("elapsed", time.Since(begin)), zap.Int64("rows", rows))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Warn(sql, zap.Error(err))
	} else {
		l.Info(sql)
	}
}
