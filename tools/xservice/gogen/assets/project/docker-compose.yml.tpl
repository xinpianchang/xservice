version: "3"
services:
  # you may change ports and environment configuration
  # run
  # docker-componse up -d --build
  {{.Name}}:
    build: .
    hostname: {{.Name}}
    ports:
      - "5000:5000"
    environment:
      - XSERVICE_ETCD=10.25.0.23:2379
      - XSERVICE_ETCD_USER=root
      - XSERVICE_ETCD_PASSWORD=123456
    restart: "unless-stopped"
