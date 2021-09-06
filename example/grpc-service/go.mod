module grpc-service

go 1.16

require (
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/labstack/echo/v4 v4.5.0
	github.com/stretchr/testify v1.7.0
	github.com/xinpianchang/xservice v1.0.20
	go.uber.org/zap v1.19.0
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/xinpianchang/xservice => ../../
