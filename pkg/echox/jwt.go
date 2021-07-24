package echox

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWT(key []byte, method *jwt.SigningMethodHMAC, skipper func(echo.Context) bool) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    key,
		SigningMethod: method.Name,
		Skipper:       skipper,
	})
}
