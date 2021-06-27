# simple configuration file, more config see `config-example.yaml`
# you may put this file to etcd (which we recommend) via following command (note, you may change --user param)
# then delete `config.yaml`, because of the configure always reads `config.yaml` in first order.
#
#  `cat config.yaml | etcdctl --user 'root:123456' put xservice/config/{{.Name}}.yaml`
#
# also, you needs to set etcd environments (user & password are optional, depends your etcd server auth configuration)
#
# export XSERVICE_ETCD=127.0.0.1:2379
# export XSERVICE_ETCD_USER=root
# export XSERVICE_ETCD_PASSWORD=123456

http:
  address: 0.0.0.0:5000

jaeger:
  agent_host: 127.0.0.1
  agent_port: 6831

log:
  level: debug
  format: console
