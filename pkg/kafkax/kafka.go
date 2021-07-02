package kafkax

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/core"
	"github.com/xinpianchang/xservice/pkg/log"
)

var (
	clients map[string]*clientWrapper
)

type config struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Broker  []string `yaml:"broker"`
}

type clientWrapper struct {
	client    sarama.Client
	name      string
	producer  sarama.SyncProducer
	countVec  *prometheus.CounterVec
	histogram prometheus.Histogram
}

func Config(v *viper.Viper) {
	var configs []config
	if err := v.UnmarshalKey("mq", &configs); err != nil {
		log.Fatal("unmarshal kafka config", zap.Error(err))
	}

	serviceName := os.Getenv(core.EnvServiceName)

	clients = make(map[string]*clientWrapper, len(configs))

	hostname, _ := os.Hostname()
	clientId := fmt.Sprint(serviceName, "_", hostname, "_", os.Getpid())

	for _, c := range configs {
		config := sarama.NewConfig()
		if version, err := sarama.ParseKafkaVersion(c.Version); err != nil {
			log.Fatal("kafka config", zap.Error(err))
		} else {
			config.Version = version
		}

		config.Consumer.Return.Errors = true
		config.Producer.Return.Successes = true
		config.ClientID = clientId
		var (
			client   sarama.Client
			producer sarama.SyncProducer
			err      error
		)
		if client, err = sarama.NewClient(c.Broker, config); err != nil {
			log.Fatal("init kafka client error", zap.Error(err), zap.String("name", c.Name))
		}

		if producer, err = sarama.NewSyncProducerFromClient(client); err != nil {
			log.Fatal("init kafka producer error", zap.Error(err), zap.String("name", c.Name))
		}

		subsystem := fmt.Sprintf("kafka_%s", strings.ReplaceAll(c.Name, "-", "_"))
		countVec := promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      "send_total",
			Help:      "Number of kafka message sent",
		}, []string{"topic"})

		histogram := promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: serviceName,
			Subsystem: subsystem,
			Name:      "send_duration_millisecond",
			Help:      "Send duration",
			Buckets:   []float64{20, 50, 100, 200, 300, 500, 1000, 2000, 3000, 5000},
		})

		clients[c.Name] = &clientWrapper{
			name:      c.Name,
			client:    client,
			producer:  producer,
			countVec:  countVec,
			histogram: histogram,
		}
	}
}

func GetClient(name ...string) sarama.Client {
	wrap := getClientWrap(name...)
	if wrap != nil {
		return wrap.client
	}
	return nil
}

func GetProducer(name ...string) sarama.SyncProducer {
	wrap := getClientWrap(name...)
	if wrap != nil {
		return wrap.producer
	}
	return nil
}

func getClientWrap(name ...string) *clientWrapper {
	var wrap *clientWrapper
	if name != nil {
		wrap = clients[name[0]]
	} else {
		wrap = clients["default"]
	}
	return wrap
}

func SendMessage(ctx context.Context, message *sarama.ProducerMessage, clientName ...string) (err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		ctx = opentracing.ContextWithSpan(context.Background(), span)
		span, _ = opentracing.StartSpanFromContext(ctx, "kafka_send")
		defer func() {
			if err != nil {
				ext.Error.Set(span, true)
				span.LogKV("err", err, "clientName", clientName)
			}
			span.Finish()
		}()
	}

	producerNotFoundErr := errors.Errorf("producer not found, clientName:%v", clientName)

	cw := getClientWrap(clientName...)
	if cw == nil {
		err = producerNotFoundErr
		return
	}
	producer := cw.producer
	if producer == nil {
		err = producerNotFoundErr
		return
	}
	start := time.Now()
	cw.countVec.WithLabelValues(message.Topic).Inc()
	_, _, err = producer.SendMessage(message)
	duration := time.Since(start).Milliseconds()
	cw.histogram.Observe(float64(duration))
	if err != nil {
		sentry.CaptureException(errors.WithMessage(err, fmt.Sprint("kafka send message to topic:", message.Topic)))
	}
	return
}

func StartGroupConsume(group string, topics []string, handler sarama.ConsumerGroupHandler, name ...string) {
	l := log.Named("kafka consumer").With(zap.String("groupID", group), zap.Strings("topics", topics))

	if group == "" || len(topics) == 0 {
		return
	}

	client := GetClient(name...)
	if client == nil {
		log.Fatal("kafkaClient not found", zap.Any("name", name))
		return
	}

	consumer, err := sarama.NewConsumerGroupFromClient(group, client)
	if err != nil {
		l.Fatal("init consume group error", zap.Error(err))
		return
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

		ctx := context.Background()
		for {
			l.Debug("start consume client ...")
			err = consumer.Consume(ctx, topics, handler)
			if err != nil {
				l.Warn("consume", zap.Error(err))
				time.Sleep(time.Second * 5)
			}
		}
	}()
}
