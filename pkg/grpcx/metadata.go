package grpcx

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

// GetIncomingMetaData get grpc incoming metadata from grpc context
func GetIncomingMetaData(ctx context.Context) (metadata.MD, bool) {
	return metadata.FromIncomingContext(ctx)
}

// GetMetaDataFirst get metadata first value by key
func GetMetaDataFirst(md metadata.MD, key string) string {
	if md == nil {
		return ""
	}

	if v := md.Get(key); len(v) > 0 {
		return v[0]
	}
	return ""
}

// GetUserAgent get user agent from grpc context
func GetUserAgent(ctx context.Context) string {
	md, _ := GetIncomingMetaData(ctx)
	return GetMetaDataFirst(md, "user-agent")
}

// GetRealIP get user ip from grpc context
func GetRealIP(ctx context.Context) string {
	md, _ := GetIncomingMetaData(ctx)
	xff := GetMetaDataFirst(md, "x-forwarded-for")
	if i := strings.IndexAny(xff, ","); i > 0 {
		return strings.TrimSpace(xff[:i])
	}

	if xrip := GetMetaDataFirst(md, "x-real-ip"); len(xrip) > 0 {
		return xrip
	}

	return ""
}
