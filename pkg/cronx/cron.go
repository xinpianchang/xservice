package cronx

import (
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/signalx"
)

const (
	SpecManual              = "@manual"
	SpecMonthly             = "@monthly"
	SpecWeekly              = "@weekly"
	SpecHourly              = "@hourly"
	SpecDaily               = "@daily"
	SpecMidnight            = "@midnight"
	SpecEveryMinutes        = "@every 1m"
	SpecEveryTenMinutes     = "@every 10m"
	SpecEveryFifteenMinutes = "@every 15m"
	SpecEveryThirtyMinutes  = "@every 30m"
	SpecEverySeconds        = "@every 1s"
	SpecEveryTenSeconds     = "@every 10s"
)

var (
	c    = cron.New(cron.WithSeconds())
	once sync.Once
)

func start() {
	once.Do(func() {
		signalx.AddShutdownHook(func(os.Signal) {
			log.Info("shutdown cron")
			c.Stop()
		})
		go c.Start()
	})
}

func Add(name string, spec string, fn func()) {
	var (
		id  cron.EntryID
		err error
	)

	getEntry := func() cron.Entry {
		return c.Entry(id)
	}

	id, err = c.AddFunc(spec, func() {
		l := log.Named(name)
		start := time.Now()
		defer func() {
			if x := recover(); x != nil {
				l.Error("panic", zap.Any("err", x))
				return
			}
			entry := getEntry()
			l.Info("done", zap.Duration("escape", time.Since(start)), zap.Time("next", entry.Next))
		}()
		fn()
	})

	if err != nil {
		log.Error("cron add", zap.Error(err))
	}

	start()
}
