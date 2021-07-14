package xservice

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/config"
	"github.com/xinpianchang/xservice/pkg/netx"
)

// Options for xservice core option
type Options struct {
	Name               string
	Version            string
	Build              string
	Description        string
	Config             *viper.Viper
	GrpcOptions        []grpc.ServerOption
	SentryOptions      sentry.ClientOptions
	EchoTracingSkipper middleware.Skipper
}

// Option for option config
type Option func(*Options)

// Name set service name
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Version set service version
func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// Build set build version for service
func Build(build string) Option {
	return func(o *Options) {
		o.Build = build
	}
}

// Description set description for service
func Description(description string) Option {
	return func(o *Options) {
		o.Description = description
	}
}

// Config set custom viper instance for xservice configuration
//
// Note: default configuration enabled watch feature, if set custom viper,
// watch feature should implements by you self
func Config(config *viper.Viper) Option {
	return func(o *Options) {
		o.Config = config
	}
}

// WithGrpcOptions add additional grpc server option
func WithGrpcOptions(options ...grpc.ServerOption) Option {
	return func(o *Options) {
		o.GrpcOptions = options
	}
}

// WithSentry set sentry option for enable sentry
func WithSentry(options sentry.ClientOptions) Option {
	return func(o *Options) {
		o.SentryOptions = options
	}
}

// WithEchoTracingSkipper set tracing skipper
func WithEchoTracingSkipper(skipper middleware.Skipper) Option {
	return func(o *Options) {
		o.EchoTracingSkipper = skipper
	}
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	loadEnvOptions(opts)

	for _, option := range options {
		option(opts)
	}

	if opts.Name == "" {
		opts.Name = core.DefaultServiceName
	}

	nameexp := `^[a-zA-Z0-9\-\_\.]+$`
	if ok, _ := regexp.MatchString(nameexp, opts.Name); !ok {
		log.Fatal("invalid service name", zap.String("name", opts.Name), zap.String("suggest", nameexp))
	}
	os.Setenv(core.EnvServiceName, opts.Name)

	if opts.Version == "" {
		opts.Version = "v0.0.1"
	}
	os.Setenv(core.EnvServiceVersion, opts.Version)

	if opts.Build == "" {
		opts.Build = fmt.Sprint("dev-", time.Now().UnixNano())
	}

	if opts.Config == nil {
		opts.loadConfig()
	}

	// env addvice addr, high priority
	if envAdvertisedAddr := os.Getenv(core.EnvAdvertisedAddr); envAdvertisedAddr != "" {
		opts.Config.SetDefault(core.ConfigServiceAdvertisedAddr, envAdvertisedAddr)
	}

	addviceAddr := opts.Config.GetString(core.ConfigServiceAdvertisedAddr)
	if addviceAddr == "" {
		address := opts.Config.GetString(core.ConfigServiceAddr)
		_, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Fatal("invalid address", zap.Error(err))
		}
		addviceAddr = net.JoinHostPort(netx.InternalIp(), port)

		opts.Config.SetDefault(core.ConfigServiceAdvertisedAddr, addviceAddr)
	}

	return opts
}

func loadEnvOptions(opts *Options) {
	if opts.Name == "" {
		opts.Name = os.Getenv(core.EnvServiceName)
	}

	if opts.Version == "" {
		opts.Version = os.Getenv(core.EnvServiceVersion)
	}
}

func (t *Options) loadConfig() {
	if err := config.LoadGlobal(); err != nil {
		log.Fatal("load config", zap.Error(err))
	}
	t.Config = viper.GetViper()
}
