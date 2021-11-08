common_protos := \
    google/api/annotations.proto \
		google/api/httpbody.proto \
		google/api/http.proto \
		google/protobuf/empty.proto \
		google/protobuf/any.proto \
		google/protobuf/struct.proto \
		validate/validate.proto \
    protoc-gen-openapiv2/options/annotations.proto \
    protoc-gen-openapiv2/options/openapiv2.proto

## all: all
.PHONY: all
all: clean build

## build: generate
.PHONY: build
build:
	@buf generate

## clean: clean build
.PHONY: clean
clean:
	@find gen ! -name 'swagger.go' ! -path gen | xargs rm -rf

.PHONY: help
help: Makefile
	@echo " Choose a command run in $(CMD):"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
