package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func Prometheus(namespace, subsystem string) echo.MiddlewareFunc {
	requests := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_total",
		Help:      "Number of requests",
	}, []string{"status", "method", "handler"})

	durations := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_duration_millisecond",
		Help:      "Request duration",
		Buckets:   []float64{50, 100, 200, 300, 500, 1000, 2000, 3000, 5000},
	}, []string{"method", "handler"})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			path := c.Request().URL.Path
			start := time.Now()
			err := next(c)
			durations.WithLabelValues(method, path).Observe(float64(time.Since(start).Milliseconds()))
			requests.WithLabelValues(fmt.Sprint(c.Response().Status), method, path).Inc()
			return err
		}
	}
}
