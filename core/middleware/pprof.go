package middleware

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

// Pprof is pprof middleware
func Pprof() echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(handler http.Handler) http.Handler {
		return http.DefaultServeMux
	})
}
