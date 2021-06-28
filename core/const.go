package core

type ContextKey string

const (
	// context
	ContextHeaderXRequestID ContextKey = "X-Request-ID" // requestId key

	DefaultServiceName = "xservice-default"

	// env key
	EnvServiceName    = "XSERVICE_NAME"
	EnvServiceVersion = "XSERVICE_VERSION"
	EnvEtcd           = "XSERVICE_ETCD"
	EnvEtcdUser       = "XSERVICE_ETCD_USER"
	EnvEtcdPassword   = "XSERVICE_ETCD_PASSWORD"

	// config key
	ConfigServiceAddr           = "http.address"
	ConfigServiceAdvertisedAddr = "http.advertised_address"

	ServiceConfigKeyPrefix   = "xservice/config"
	ServiceRegisterKeyPrefix = "xservice/register"
)
