package xservice

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/google/gops/agent"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/gormx"
	"github.com/xinpianchang/xservice/pkg/kafkax"
	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/redisx"
	"github.com/xinpianchang/xservice/pkg/tracingx"
)

// Service is a service interface, core api for xservice
type Service interface {
	Name() string
	Options() *Options
	Client() Client
	Server() Server
	String() string
}

type serviceImpl struct {
	options *Options
	client  Client
	server  Server
}

// New create new xservice instance
func New(options ...Option) Service {
	opts := loadOptions(options...)

	service := new(serviceImpl)
	service.options = opts

	service.init()

	return service
}

// Name get service name
func (t *serviceImpl) Name() string {
	return t.options.Name
}

// Options get service options
func (t *serviceImpl) Options() *Options {
	return t.options
}

// Client get service client
func (t *serviceImpl) Client() Client {
	return t.client
}

// Server get service server
func (t *serviceImpl) Server() Server {
	return t.server
}

// String get service formatted name
func (t *serviceImpl) String() string {
	return fmt.Sprint(t.Name(), "/", t.options.Version, " - ", t.options.Description)
}

func (t *serviceImpl) init() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal("agent", zap.Error(err))
	}

	if t.options.Config.IsSet("log") {
		log.Config(t.options.Config)
	}

	tracingx.Config(t.options.Config)

	if t.options.SentryOptions.Dsn != "" {
		err := sentry.Init(t.options.SentryOptions)
		if err != nil {
			log.Fatal("init sentry", zap.Error(err))
		}
	}

	if t.options.Config.IsSet("redis") {
		redisx.Config(t.options.Config)
	}

	if t.options.Config.IsSet("database") {
		gormx.Config(t.options.Config)
	}

	if t.options.Config.IsSet("mq") {
		kafkax.Config(t.options.Config)
	}

	if t.options.Config.IsSet(core.ConfigServiceAddr) {
		t.server = newServer(t.options)
	}

	t.client = newClient(t.options)
}
