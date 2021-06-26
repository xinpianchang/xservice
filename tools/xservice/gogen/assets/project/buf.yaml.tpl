version: v1beta1
name: buf.build/{{.Repo}}/{{.Name}}
deps:
  - buf.build/beta/googleapis
  - buf.build/beta/protoc-gen-validate
build:
  roots:
    - .
lint:
  use:
    - DEFAULT
breaking:
  use:
    - FILE
