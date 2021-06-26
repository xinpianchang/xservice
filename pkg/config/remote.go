package config

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
)

type remoteConfig struct {
	viper.RemoteProvider
}

func (t *remoteConfig) Get(rp viper.RemoteProvider) (io.Reader, error) {
	t.RemoteProvider = rp
	return t.get()
}

func (t *remoteConfig) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	t.RemoteProvider = rp
	return t.get()
}

func (t *remoteConfig) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	t.RemoteProvider = rp

	rr := make(chan *viper.RemoteResponse)
	stop := make(chan bool)

	go func() {
		client, err := t.client()
		if err != nil {
			log.Fatal("watch config channel", zap.Error(err))
			return
		}

		defer client.Close()

		for {
			ch := client.Watch(context.Background(), t.Path())

			select {
			case <-stop:
				return
			case res := <-ch:
				for _, event := range res.Events {
					rr <- &viper.RemoteResponse{
						Value: event.Kv.Value,
					}
				}
			}
		}
	}()

	return rr, stop
}

func (t *remoteConfig) client() (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   []string{t.Endpoint()},
		DialTimeout: time.Second * 5,
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

func (t *remoteConfig) get() (io.Reader, error) {
	client, err := t.client()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	resp, err := client.Get(context.Background(), t.Path())
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) > 0 {
		return bytes.NewReader(resp.Kvs[0].Value), nil
	}
	return strings.NewReader(""), nil
}
