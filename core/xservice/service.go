package xservice

import (
	"fmt"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
)

func serviceKey(serviceName string, desc *grpc.ServiceDesc) string {
	host, _ := os.Hostname()
	if host == "" {
		host = fmt.Sprint("unknown-host-pid-", os.Getpid())
	}
	return fmt.Sprint(serviceKeyPrefix(serviceName, desc), "/", host)
}

func serviceKeyPrefix(serviceName string, desc *grpc.ServiceDesc) string {
	return fmt.Sprint(core.ServiceRegisterKeyPrefix, "/", serviceName, "/", desc.ServiceName)
}

func serviceEtcdClient() (*clientv3.Client, error) {
	endpoints := strings.Split(os.Getenv(core.EnvEtcd), ",")

	cfg := clientv3.Config{
		Endpoints:         endpoints,
		DialTimeout:       time.Second * 5,
		DialKeepAliveTime: time.Second * 5,
		AutoSyncInterval:  time.Second * 10,
		Logger:            log.Get().WithOptions(zap.IncreaseLevel(zapcore.ErrorLevel)),
	}
	if username := os.Getenv(core.EnvEtcdUser); username != "" {
		cfg.Username = username
	}
	if password := os.Getenv(core.EnvEtcdPassword); password != "" {
		cfg.Password = password
	}
	client, err := clientv3.New(cfg)
	return client, err
}
