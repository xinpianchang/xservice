package cronx

import (
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/v2/pkg/log"
	"github.com/xinpianchang/xservice/v2/pkg/signalx"
)

const (
	SpecManual              = "@manual"    // specification for manual
	SpecMonthly             = "@monthly"   // spec monthly
	SpecWeekly              = "@weekly"    // spec weekly
	SpecHourly              = "@hourly"    // spec hourly
	SpecDaily               = "@daily"     // spec daily
	SpecMidnight            = "@midnight"  // spec at midnight
	SpecEveryMinutes        = "@every 1m"  // spec every 1 minute
	SpecEveryTenMinutes     = "@every 10m" // spec every 10 minutes
	SpecEveryFifteenMinutes = "@every 15m" // spec every 15 minutes
	SpecEveryThirtyMinutes  = "@every 30m" // spec every 30 minutes
	SpecEverySeconds        = "@every 1s"  // spec every 1 second
	SpecEveryTenSeconds     = "@every 10s" // spec every 10 seconds
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

	entry := getEntry()
	log.Debug("cron add", zap.String("name", name), zap.String("spec", spec), zap.Time("next", entry.Schedule.Next(time.Now())))

	start()
}
