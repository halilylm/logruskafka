package logruskafka

import (
	"context"
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

// KafkaAuth when sasl needed
// don't fill algorithm when you only
// plain auth, fill it when you
// need scram
type KafkaAuth struct {
	// Username for auth
	Username string
	// Password for auth
	Password string
	// Algorithm to be used in scram
	Algorithm scram.Algorithm
	// TLSConfig for transportation
	TLSConfig tls.Config
}

type KafkaHook struct {
	// id of the hook
	id string
	// levels allowed
	levels []logrus.Level
	// formatter for log entry
	formatter logrus.Formatter
	// writer for kafka
	w *kafka.Writer
}

// NewKafkaHook generates a new kafka hook
func NewKafkaHook(
	id string,
	levels []logrus.Level,
	formatter logrus.Formatter,
	topic string,
	brokers []string,
) (*KafkaHook, error) {
	writer := newWriterWithoutAuth(topic, brokers)
	return &KafkaHook{
		id:        id,
		levels:    levels,
		formatter: formatter,
		w:         writer,
	}, nil
}

// NewKafkaHookWithSaslAuth generates a new kafka hook
// with sasl auth
func NewKafkaHookWithSaslAuth(
	id string,
	levels []logrus.Level,
	formatter logrus.Formatter,
	topic string,
	brokers []string,
	auth *KafkaAuth) (*KafkaHook, error) {
	writer, err := newWriterWithAuth(topic, brokers, auth)
	if err != nil {
		return nil, err
	}
	return &KafkaHook{
		id:        id,
		levels:    levels,
		formatter: formatter,
		w:         writer,
	}, nil
}

// ID returns id of the hook
func (k *KafkaHook) ID() string {
	return k.id
}

// Levels return logging levels
func (k *KafkaHook) Levels() []logrus.Level {
	return k.levels
}

// Fire will process the log
func (k *KafkaHook) Fire(entry *logrus.Entry) error {
	// Format before writing
	b, err := k.formatter.Format(entry)
	if err != nil {
		return err
	}

	if err := k.w.WriteMessages(context.TODO(), kafka.Message{Value: b}); err != nil {
		return err
	}
	return nil
}

// newWriterWithoutAuth generate a new writer
// when no auth required
func newWriterWithoutAuth(topic string, brokers []string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
		Topic:    topic,
	}
}

// newWriterWithAuth generate a new writer
// when auth required
func newWriterWithAuth(topic string, brokers []string, auth *KafkaAuth) (*kafka.Writer, error) {
	var mechanism sasl.Mechanism
	var err error
	// if no algorithm provided
	// we will assume it is plain
	// mechanism
	if auth.Algorithm == nil {
		mechanism = plain.Mechanism{
			Username: "username",
			Password: "password",
		}
	} else {
		mechanism, err = scram.Mechanism(auth.Algorithm, auth.Username, auth.Password)
		if err != nil {
			return nil, err
		}
	}
	sharedTransport := &kafka.Transport{
		SASL: mechanism,
		TLS:  &auth.TLSConfig,
	}

	w := &kafka.Writer{
		Addr:      kafka.TCP(brokers...),
		Topic:     topic,
		Balancer:  &kafka.Hash{},
		Transport: sharedTransport,
	}
	return w, nil
}
