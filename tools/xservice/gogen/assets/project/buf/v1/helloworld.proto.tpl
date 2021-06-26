syntax = "proto3";

package buf.v1;
option go_package = "{{.Module}}/buf/v1";

import "google/api/httpbody.proto";
import "google/api/annotations.proto";
import "google/protobuf/empty.proto";
import "google/protobuf/any.proto";
import "google/protobuf/struct.proto";
import "validate/validate.proto";

// HelloWorldService
service HelloWorldService {
  // Hello
  //
  // {{`{{import "buf/v1/tables.md"}}`}}
  rpc Hello(HelloRequest) returns (HelloResponse) {
    option(google.api.http) = {
      post: "/rpc/v1/hello"
      body: "*"
    };
  }
}

message HelloRequest {
  string name = 1 [(validate.rules).string = {min_len: 1, max_len: 32}];
}

message HelloResponse {
  string message = 1;
}
