FROM golang:alpine AS builder
RUN apk update --no-cache && apk add git
LABEL stage=gobuilder
WORKDIR /build/xservice
COPY . .
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct
RUN go run mage.go -v build

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai
WORKDIR /app
COPY --from=builder /build/xservice/dist/{{.Name}}-linux-amd64/ /app/
CMD ["./{{.Name}}"]
