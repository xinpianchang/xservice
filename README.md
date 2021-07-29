# xservice [WIP]

[![License](https://img.shields.io/github/license/xinpianchang/xservice?style=flat-square)](https://raw.githubusercontent.com/xinpianchang/xservice/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/xinpianchang/xservice?style=flat-square)](https://goreportcard.com/report/github.com/xinpianchang/xservice)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/xinpianchang/xservice)

Another excellent micro service framework

## Features

- RESTful API (base on echo/v4)
- gRPC & gRPC gateway service & Swagger document generation
- Service discovery (base on ETCD/v3)
- gRPC & gRPC-Gateway & RESTful API all in one tcp port, mux via `cmux`
- Builtin middlewares & easily to extended
- Prometheus & Tracing (jaeger) & Sentry integrated
- Embed toolset for code generation (e.g. GORM & model CRUD & project layout)

## Quick start

Install toolset.

```bash
go install github.com/xinpianchang/xservice/tools/xservice@latest
```

Create new project via toolset.

```bash
mkdir hello
cd hello
xservice new --module github.com/example/hello
```

Open the generated `README.md` file, following the initialize steps and happing coding. 🎉

## Resource

- go-zero https://github.com/tal-tech/go-zero (special thanks)
- micro https://github.com/asim/go-micro (inspired)
- GORM http://gorm.io/
- Echo https://echo.labstack.com/
- validator https://github.com/go-playground/validator
- gRPC generate tool/buf https://buf.build/
- gRPC validate https://github.com/envoyproxy/protoc-gen-validate
- RESTful validate https://github.com/go-playground/validator
- gRPC-Gateway https://grpc-ecosystem.github.io/grpc-gateway/
- jaeger https://www.jaegertracing.io/
