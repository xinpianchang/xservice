package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/xinpianchang/xservice/core/xservice"
	"github.com/xinpianchang/xservice/pkg/swaggerui"

	pb "{{.Module}}/buf/v1"
	"{{.Module}}/service"
	"{{.Module}}/version"
)

var (
  showVersion = flag.Bool("version", false, "print version")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s version:%s, build:%s, runtime:%s\n",
			version.Name, version.Version, version.Build, runtime.Version())
		return
	}

	srv := xservice.New(
		xservice.Name(version.Name),
		xservice.Version(version.Version),
		xservice.Build(version.Build),
		xservice.Description(version.Description),
	)

	server := srv.Server()

	// swagger doc
	server.Echo().Group("/swagger/*", swaggerui.Serve("/swagger/", pb.SwaggerFS))

	// register grpc service
	server.GrpcRegister(&pb.HelloWorldService_ServiceDesc, &service.HelloWorldServiceServerImpl{}, pb.RegisterHelloWorldServiceHandler)

	// routes config
	routes(server.Echo())

	if err := server.Serve(); err != nil {
		panic(err)
	}
}

// routes for RESTful api
func routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
}
