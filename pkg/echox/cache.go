package echox

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func Cache(age time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int64(age.Seconds())))
			return next(c)
		}
	}
}

func NoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
		return next(c)
	}
}
