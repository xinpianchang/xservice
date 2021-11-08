version: v1
deps:
  - buf.build/googleapis/googleapis
  - buf.build/envoyproxy/protoc-gen-validate
lint:
  use:
    - DEFAULT
breaking:
  use:
    - FILE
