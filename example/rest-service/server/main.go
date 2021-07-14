package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xinpianchang/xservice/core/xservice"
)

func main() {
	srv := xservice.New(
		xservice.Name("rest-service"),
		xservice.Version("v1.0.0"),
		xservice.Description("example RESTFul service"),
	)

	routes(srv.Server().Echo())

	srv.Server().Serve()
}

func routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
}
