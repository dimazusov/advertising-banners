package kafka

import (
	"context"
	"time"

	"github.com/dimazusov/hw-test/advertising-banners/internal/config"
	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
}

type Consumer interface {
	Listen(f func(b []byte) error) error
}

func NewConsumer(topic string, cfg *config.Config) (Consumer, error) {
	dialerTimeout, err := time.ParseDuration(cfg.Kafka.DialerTimeout)
	if err != nil {
		return nil, err
	}

	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{cfg.Kafka.Address},
			Topic:     topic,
			GroupID:   cfg.Kafka.Group,
			Partition: 0,
			MinBytes:  10e3, // 10KB
			MaxBytes:  50e6, // 50MB
			Dialer: &kafka.Dialer{
				Timeout:   dialerTimeout,
				DualStack: true,
			},
		},
		)}, nil
}

func (m *consumer) Listen(f func(b []byte) error) error {
	defer m.reader.Close()

	for {
		message, err := m.reader.FetchMessage(context.Background())
		if err != nil {
			return err
		}

		err = f(message.Value)
		if err != nil {
			return err
		}

		err = m.reader.CommitMessages(context.Background(), message)
		if err != nil {
			return err
		}
	}
}
