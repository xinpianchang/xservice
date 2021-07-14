package core

// ContextKey for context key, string alias
type ContextKey string

const (
	// context
	ContextHeaderXRequestID ContextKey = "X-Request-ID" // requestId key

	DefaultServiceName = "xservice-default" // default service name

	// env key
	EnvServiceName    = "XSERVICE_NAME"            // serviceName key
	EnvServiceVersion = "XSERVICE_VERSION"         // serviceVersion key
	EnvAdvertisedAddr = "XSERVICE_ADVERTISED_ADDR" // advertised addr key
	EnvEtcd           = "XSERVICE_ETCD"            // etc endpoint key
	EnvEtcdUser       = "XSERVICE_ETCD_USER"       // etcdUser key
	EnvEtcdPassword   = "XSERVICE_ETCD_PASSWORD"   // etcdPassword key

	// config key
	ConfigServiceAddr           = "http.address"            // config http address key
	ConfigServiceAdvertisedAddr = "http.advertised_address" // config http advertise address key

	ServiceConfigKeyPrefix   = "xservice/config"   // service config key prefix
	ServiceRegisterKeyPrefix = "xservice/register" // service register key prefix
)
