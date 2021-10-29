version: v1beta1
name: buf.build/{{.Repo}}/{{.Name}}
deps:
  - buf.build/googleapis/googleapis
  - buf.build/envoyproxy/protoc-gen-validate
build:
  roots:
    - buf
lint:
  use:
    - DEFAULT
breaking:
  use:
    - FILE
