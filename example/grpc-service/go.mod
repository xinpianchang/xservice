module grpc-service

go 1.16

require (
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/labstack/echo/v4 v4.3.0
	github.com/stretchr/testify v1.7.0
	github.com/xinpianchang/xservice v1.0.0
	go.uber.org/zap v1.17.0
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/xinpianchang/xservice => ../../
