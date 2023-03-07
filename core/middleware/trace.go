package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xinpianchang/xservice/v2/core"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

// Trace is opentelemetry middleware
func Trace(skipper middleware.Skipper) echo.MiddlewareFunc {
	serviceName := os.Getenv(core.EnvServiceName)
	return otelecho.Middleware(
		serviceName,
		otelecho.WithPropagators(otel.GetTextMapPropagator()),
		otelecho.WithSkipper(skipper),
	)
}
