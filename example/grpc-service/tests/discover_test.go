package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xinpianchang/xservice/v2/core/xservice"

	pb "grpc-service/buf/v1"
)

func Test_service_discover(t *testing.T) {
	srv := xservice.New(
		xservice.Name("grpc-service"),
	)

	ctx := context.Background()

	conn, err := srv.Client().GrpcClientConn(ctx, "grpc-service", &pb.CalculatorService_ServiceDesc)
	require.NoError(t, err)
	require.NotNil(t, conn)

	client := pb.NewCalculatorServiceClient(conn)
	response, err := client.AddInt(ctx, &pb.AddIntRequest{A: 1, B: 2})
	require.NoError(t, err)
	require.Equal(t, int32(3), response.Result)
}
