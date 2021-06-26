package echox

import (
	"github.com/dgrijalva/jwt-go"
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
