module grpc-service

go 1.16

require (
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3
	github.com/labstack/echo/v4 v4.7.2
	github.com/stretchr/testify v1.7.5
	github.com/xinpianchang/xservice v1.0.20
	go.uber.org/zap v1.21.0
	google.golang.org/genproto v0.0.0-20220627200112-0a929928cb33
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/xinpianchang/xservice => ../../
