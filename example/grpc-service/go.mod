module grpc-service

go 1.16

require (
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/labstack/echo/v4 v4.7.2
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/prometheus/common v0.34.0 // indirect
	github.com/stretchr/testify v1.7.2
	github.com/subosito/gotenv v1.4.0 // indirect
	github.com/xinpianchang/xservice v1.0.20
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/ini.v1 v1.66.6 // indirect
)

replace github.com/xinpianchang/xservice => ../../
