package kafkax

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/netx"
)

var (
	configMap = make(map[string]mqConfig, 8)
	clientID  int32
)

type mqConfig struct {
	Name    string   `yaml:"name"`    // config name, should be unique
	Version string   `yaml:"version"` // kafka cluster version
	Broker  []string `yaml:"broker"`  // kafka broker list
}

func Config(v *viper.Viper) {
	var configs []mqConfig
	if err := v.UnmarshalKey("mq", &configs); err != nil {
		log.Fatal("unmarshal kafka config", zap.Error(err))
	}

	for _, it := range configs {
		configMap[it.Name] = it
	}
}

// Client is a kafka client wrapper
type Client interface {
	// Get get kafka client
	Get() sarama.Client

	// Close close kafka client
	Close() error

	// SendMessage send message to kafka
	SendMessage(ctx context.Context, message *sarama.ProducerMessage) error

	// GroupConsume consume kafka group
	GroupConsume(ctx context.Context, group string, topics []string, handler sarama.ConsumerGroupHandler) error
}

type defaultKafka struct {
	config                 mqConfig
	client                 sarama.Client
	producerInitializeOnce sync.Once
	producer               sarama.SyncProducer
}

// New create a kafka client
func New(name string, cfg ...*sarama.Config) (Client, error) {
	c, ok := configMap[name]
	if !ok {
		return nil, errors.Errorf("configuration not found, name: %v", name)
	}

	kafka := &defaultKafka{
		config: c,
	}

	var config *sarama.Config
	if len(cfg) > 0 {
		config = cfg[0]
	}
	if config == nil {
		config = NewDefaultKafkaConfig()
	}

	if version, err := sarama.ParseKafkaVersion(c.Version); err != nil {
		log.Fatal("kafka config", zap.Error(err))
	} else {
		config.Version = version
	}

	if config.ClientID == "" {
		config.ClientID = kafkaClientID()
	}

	if client, err := sarama.NewClient(c.Broker, config); err != nil {
		log.Fatal("init kafka client error", zap.Error(err), zap.String("name", c.Name))
	} else {
		kafka.client = client
	}

	return kafka, nil
}

func NewDefaultKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	maxMessageBytes := 1024 * 1024 * 10
	if config.Producer.MaxMessageBytes < maxMessageBytes {
		config.Producer.MaxMessageBytes = maxMessageBytes
	}

	config.ClientID = kafkaClientID()

	return config
}

func kafkaClientID() string {
	atomic.AddInt32(&clientID, 1)
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = netx.InternalIp()
	}
	return fmt.Sprint("sarama", "_", hostname, "_", os.Getpid(), "_", clientID)
}

// Get get kafka client
func (t *defaultKafka) Get() sarama.Client {
	return t.client
}

// Close close kafka client
func (t *defaultKafka) Close() error {
	if !t.client.Closed() {
		return t.client.Close()
	}
	return nil
}

// SendMessage send message to kafka
func (t *defaultKafka) SendMessage(ctx context.Context, message *sarama.ProducerMessage) (err error) {
	if err := t.initializeProducer(ctx); err != nil {
		return err
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		ctx = opentracing.ContextWithSpan(context.Background(), span)
		span, _ = opentracing.StartSpanFromContext(ctx, "kafka_send")
		defer func() {
			if err != nil {
				ext.Error.Set(span, true)
				span.LogKV("client_name", t.config.Name, "err", err)
			}
			span.Finish()
		}()
	}

	_, _, err = t.producer.SendMessage(message)

	return
}

// GroupConsume consume kafka group
func (t *defaultKafka) GroupConsume(ctx context.Context, group string, topics []string, handler sarama.ConsumerGroupHandler) error {
	l := log.Named("kafka consumer").With(
		zap.String("name", t.config.Name),
		zap.String("groupId", group),
		zap.Strings("topics", topics),
	)

	if group == "" || len(topics) == 0 {
		return nil
	}

	consumer, err := sarama.NewConsumerGroupFromClient(group, t.client)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if x := recover(); x != nil {
				l.Error("start consumer error", zap.Error(err))
			}
		}()

		go func() {
			for it := range consumer.Errors() {
				l.Warn("client error", zap.Error(it))
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}
			l.Debug("start consume")
			err = consumer.Consume(ctx, topics, handler)
			if err != nil {
				l.Warn("consume", zap.Error(err))
				time.Sleep(time.Second * 5)
			}
		}
	}()

	return nil
}

func (t *defaultKafka) initializeProducer(ctx context.Context) (err error) {
	t.producerInitializeOnce.Do(func() {
		t.producer, err = sarama.NewSyncProducerFromClient(t.client)
	})
	return
}
