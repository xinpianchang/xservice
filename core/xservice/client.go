package xservice

import (
	"context"
	"fmt"
	"os"

	resolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	gresolver "google.golang.org/grpc/resolver"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
)

type Client interface {
	GrpcClientConn(ctx context.Context, service string, desc *grpc.ServiceDesc, endpoint ...string) (*grpc.ClientConn, error)
}

type clientImpl struct {
	options  *Options
	resolver gresolver.Builder
}

func newClient(opts *Options) Client {
	client := &clientImpl{
		options: opts,
	}

	if os.Getenv(core.EnvEtcd) != "" {
		cli, err := serviceEtcdClient()
		if err != nil {
			log.Fatal("etcd client", zap.Error(err))
		}
		client.resolver, err = resolver.NewBuilder(cli)
		if err != nil {
			log.Fatal("endpoints manager", zap.Error(err))
		}
	}

	return client
}

func (t *clientImpl) GrpcClientConn(ctx context.Context, service string, desc *grpc.ServiceDesc, endpoint ...string) (*grpc.ClientConn, error) {
	if len(endpoint) > 0 {
		return grpc.DialContext(ctx, endpoint[0], grpc.WithInsecure())
	}

	if os.Getenv(core.EnvEtcd) == "" {
		log.Fatal("etcd not configured")
	}
	target := fmt.Sprint("etcd:///", serviceKeyPrefix(service, desc))
	log.For(ctx).Debug("client conn", zap.String("target", target))
	return grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithResolvers(t.resolver))
}
