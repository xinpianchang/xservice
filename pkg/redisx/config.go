package redisx

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/v2/pkg/log"
	"github.com/xinpianchang/xservice/v2/pkg/signalx"
)

var (
	clients map[string]*redis.Client
	Locker  *redislock.Client
	cfgMap  map[*redis.Client]*redisConfig
)

type redisConfig struct {
	Name         string `yaml:"name"`
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	ReadTimeout  int    `yaml:"readTimeout"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"poolSize"`
	MaxRetries   int    `yaml:"maxRetries"`
	MinIdleConns int    `yaml:"minIdleConns"`
	MaxConnAge   int    `yaml:"maxConnAge"`
	Prefix       string `yaml:"prefix"`
}

func Config(v *viper.Viper) {
	var cfg []*redisConfig

	if err := v.UnmarshalKey("redis", &cfg); err != nil {
		log.Fatal("read redis config", zap.Error(err))
	}

	clients = make(map[string]*redis.Client, len(cfg))
	cfgMap = make(map[*redis.Client]*redisConfig, len(cfg))
	for _, c := range cfg {
		if c.MinIdleConns <= 0 {
			c.MinIdleConns = 0
		}
		if c.MaxRetries < 0 {
			c.MaxRetries = 2
		}
		if c.PoolSize < 0 {
			c.PoolSize = 500
		}
		if c.ReadTimeout <= 0 {
			c.ReadTimeout = 5
		}
		if c.MaxConnAge <= 0 {
			c.MaxConnAge = 300
		}

		if c.Prefix == "" {
			data := []byte(fmt.Sprint(time.Now().UnixNano()))
			hash := md5.Sum(data)
			prefix := hex.EncodeToString(hash[:4])
			c.Prefix = prefix
		}

		client := redis.NewClient(&redis.Options{
			Addr:            c.Addr,
			Password:        c.Password,
			DB:              c.DB,
			ReadTimeout:     time.Second * time.Duration(c.ReadTimeout),
			MaxRetries:      c.MaxRetries,
			MinIdleConns:    c.MinIdleConns,
			ConnMaxLifetime: time.Second * time.Duration(c.MaxConnAge),
			PoolSize:        c.PoolSize,
		})
		r := client.Ping(context.TODO())
		if err := r.Err(); err != nil {
			log.Fatal("ping", zap.String("name", c.Name), zap.Error(err))
		}

		// log.Debug(fmt.Sprint("redis ", c.Name, " ping"),
		// 	zap.String("rsp", r.Val()), zap.String("addr", c.Addr), zap.Int("db", c.DB))

		client.AddHook(&redisTracing{})

		clients[c.Name] = client
		cfgMap[client] = c

		// as main redis
		if c.Name == "redis" && Locker == nil {
			Locker = redislock.New(client)
		}
	}

	signalx.AddShutdownHook(func(os.Signal) {
		for _, c := range clients {
			_ = c.Close()
		}
	})
}

func GetClient(name string) *redis.Client {
	return clients[name]
}

func NewLocker(client *redis.Client) *redislock.Client {
	return redislock.New(client)
}

func Key(client *redis.Client, keys ...string) string {
	if c, ok := cfgMap[client]; ok {
		keys = append(append(make([]string, 0, len(keys)+1), c.Prefix), keys...)
	}
	return strings.Join(keys, ":")
}

type ClientWrapper struct {
	*redis.Client
}

func Get(name string) *ClientWrapper {
	return &ClientWrapper{GetClient(name)}
}

func (t *ClientWrapper) Key(keys ...string) string {
	return Key(t.Client, keys...)
}

func (t *ClientWrapper) Obtain(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (*redislock.Lock, error) {
	locker := NewLocker(t.Client)
	return locker.Obtain(ctx, key, ttl, opt)
}
