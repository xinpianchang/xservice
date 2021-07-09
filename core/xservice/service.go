package xservice

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
)

// serviceKey get full service key
func serviceKey(serviceName string, desc *grpc.ServiceDesc) string {
	host, _ := os.Hostname()
	pid := os.Getpid()
	if host == "" {
		host = "unknown-host"
	}
	return fmt.Sprint(serviceKeyPrefix(serviceName, desc), "/", host, "-pid-", pid)
}

// serviceKeyPrefix get service key prefix
func serviceKeyPrefix(serviceName string, desc *grpc.ServiceDesc) string {
	return fmt.Sprint(core.ServiceRegisterKeyPrefix, "/", serviceName, "/", desc.ServiceName)
}

var (
	_etcdClient     *clientv3.Client
	_etcdClientOnce sync.Once
)

// serviceEtcdClient lazy init etcd client
func serviceEtcdClient() *clientv3.Client {
	_etcdClientOnce.Do(func() {
		endpoints := strings.Split(os.Getenv(core.EnvEtcd), ",")
		cfg := clientv3.Config{
			Endpoints:         endpoints,
			DialTimeout:       time.Second * 5,
			DialKeepAliveTime: time.Second * 10,
			AutoSyncInterval:  time.Second * 30,
			Logger:            log.Get().WithOptions(zap.IncreaseLevel(zapcore.ErrorLevel)),
		}
		if username := os.Getenv(core.EnvEtcdUser); username != "" {
			cfg.Username = username
		}
		if password := os.Getenv(core.EnvEtcdPassword); password != "" {
			cfg.Password = password
		}
		client, err := clientv3.New(cfg)
		if err != nil {
			log.Fatal("serviceEtcdClient", zap.Error(err))
		}
		_etcdClient = client
	})
	return _etcdClient
}
