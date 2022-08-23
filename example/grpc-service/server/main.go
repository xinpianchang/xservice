package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/core/xservice"
	"github.com/xinpianchang/xservice/pkg/echox"
	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/swaggerui"

	pb "grpc-service/buf/v1"
)

func main() {
	srv := xservice.New(
		xservice.Name("grpc-service"),
		xservice.Version("v1.0.0"),
		xservice.Description("example grpc service with enable grpc gateway"),
		xservice.WithGrpcServerEnableReflection(true), // optional
	)

	server := srv.Server()

	server.Echo().Use(echox.Dump())
	server.Echo().Group("/swagger/*", swaggerui.Serve("/swagger/", pb.SwaggerFS))

	server.GrpcRegister(&pb.GreeterService_ServiceDesc, &GreeterServer{}, pb.RegisterGreeterServiceHandler)
	server.GrpcRegister(&pb.CalculatorService_ServiceDesc, &CalculatorServer{}, pb.RegisterCalculatorServiceHandler)

	server.Echo().GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	if err := server.Serve(); err != nil {
		log.Fatal("serve", zap.Error(err))
	}

	// curl -v -X POST -k http://127.0.0.1:5001/rpc/v1/echo -d '{"name": "world"}'
	// curl -v -X POST -k http://127.0.0.1:5001/rpc/v1/calculator -d '{"a": 1, "b": 2}'
}

// implementations

type GreeterServer struct{}

func (t *GreeterServer) SayHello(ctx context.Context, request *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
	return &pb.SayHelloResponse{Message: fmt.Sprint("hello ", request.Name)}, nil
}

type CalculatorServer struct{}

func (t *CalculatorServer) AddInt(ctx context.Context, request *pb.AddIntRequest) (*pb.AddIntResponse, error) {
	log.Info("call addInt", zap.Any("request", request))
	return &pb.AddIntResponse{Result: request.A + request.B}, nil
}
