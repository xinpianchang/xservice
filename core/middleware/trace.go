package middleware

import (
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
)

func Trace(bodyDump bool, skipper middleware.Skipper) echo.MiddlewareFunc {
	return jaegertracing.TraceWithConfig(jaegertracing.TraceConfig{
		Tracer:     opentracing.GlobalTracer(),
		Skipper:    skipper,
		IsBodyDump: bodyDump,
	})
}
