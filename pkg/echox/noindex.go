package echox

import (
	"github.com/labstack/echo/v4"
)

func NoIndex(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive, nosnippet")
		return next(c)
	}
}
