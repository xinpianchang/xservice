package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"github.com/xinpianchang/xservice/core"
)

func Load() (*viper.Viper, error) {
	v := viper.New()
	err := load(v)
	return v, err
}

func LoadGlobal() error {
	return load(viper.GetViper())
}

// Load auto load from local yaml config or etcd if local file not exists
func load(v *viper.Viper) error {
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	executable, _ := os.Executable()
	v.AddConfigPath(filepath.Dir(executable))

	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.AddConfigPath("../../")

	if _, file, _, ok := runtime.Caller(0); ok {
		v.AddConfigPath(filepath.Join(filepath.Dir(file), "../../"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			endpoint := os.Getenv(core.EnvEtcd)
			if endpoint != "" {
				name := os.Getenv(core.EnvServiceName)
				if name == "" {
					name = core.DefaultServiceName
				}
				viper.RemoteConfig = &remoteConfig{}
				if e := v.AddRemoteProvider("etcd", endpoint, fmt.Sprint(core.ServiceConfigKeyPrefix, "/", name, ".yaml")); e != nil {
					return e
				}
				if err = v.ReadRemoteConfig(); err != nil {
					return err
				}
				v.SetDefault("dir", filepath.Dir(executable))
				_ = v.WatchRemoteConfigOnChannel()
			}
		} else {
			return err
		}
	} else {
		v.SetDefault("dir", filepath.Dir(v.ConfigFileUsed()))
		v.WatchConfig()
	}

	return nil
}
