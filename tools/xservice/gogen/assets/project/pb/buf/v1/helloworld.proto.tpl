syntax = "proto3";

package buf.v1;
option go_package = "{{.Module}}_pb/gen/v1";

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
  // The name who you want to say hello (name's length should be `1 <= len(name) <= 32`)
  string name = 1 [(validate.rules).string = {min_len: 1, max_len: 32}];
}

message HelloResponse {
  // represent hello message response
  string message = 1;
}
