# {{.Name}}

Robust service powered by `xservice`

## Prepare

install protobuf & generator plugins

```bash
# install protoc v3.17.3 (recommended)
# https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3

# or checkout latest
# https://github.com/protocolbuffers/protobuf/releases
```

```bash
# install specific version (recommended)
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.42.0
go install github.com/envoyproxy/protoc-gen-validate@v0.6.2

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.6.0 \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.6.0



# install latest (not well tested)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/envoyproxy/protoc-gen-validate@latest

go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

install buf

```bash
go install github.com/bufbuild/buf/cmd/buf@latest
```

## Initialize project

```bash
cd pb && buf mod update buf && buf generate && cd ../
go mod tidy && go mod download && go generate
```

## Run

```bash
go run cmd/main.go
```

## Resource

- xservice https://github.com/xinpianchang/xservice
- GORM https://gorm.io/docs/ & https://github.com/go-gorm/gorm
- Echo https://echo.labstack.com/
- validator https://github.com/go-playground/validator
- gRPC https://grpc.io/
- Protobuf https://developers.google.com/protocol-buffers/docs/gotutorial
- gRPC generate tool/buf https://buf.build/
- gRPC validate https://github.com/envoyproxy/protoc-gen-validate
- RESTful validate https://github.com/go-playground/validator
- gRPC-Gateway https://grpc-ecosystem.github.io/grpc-gateway/
- jaeger https://www.jaegertracing.io/
- Prometheus https://prometheus.io/
