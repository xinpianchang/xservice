package redisx

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func newViper() *viper.Viper {
	v := viper.New()
	v.SetDefault("redis", map[string]interface{}{
		"name":     "redis",
		"addr":     "127.0.0.1:6379",
		"password": "123456",
	})

	return v
}

func Test_config(t *testing.T) {
	v := newViper()
	Config(v)
	assert.NotNil(t, Get("redis"))
	assert.NotNil(t, Locker)
}

func Test_wrapper(t *testing.T) {
	v := newViper()
	Config(v)

	w := Get("redis")
	r := w.Ping(context.Background())
	assert.NoError(t, r.Err())

	ret, err := r.Result()
	assert.NoError(t, err)
	t.Log("result:", ret)
}
