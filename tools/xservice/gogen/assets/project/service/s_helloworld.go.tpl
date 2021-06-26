package service

import (
	"context"
	"fmt"

	pb "{{.Module}}/buf/v1"
)

type HelloWorldServiceServerImpl struct{}

func (t *HelloWorldServiceServerImpl) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: fmt.Sprint("Hello ", request.Name)}, nil
}
