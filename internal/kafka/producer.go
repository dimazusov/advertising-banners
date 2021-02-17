package kafka

import (
	"github.com/dimazusov/hw-test/advertising-banners/internal/config"
	"github.com/segmentio/kafka-go"

	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

type producer struct {
	writer *kafka.Writer
}

type Producer interface {
	WriteMessages(ctx context.Context, value []byte) error
	Close() error
}

func NewProducer(topic string, cfg *config.Config) (Producer, error) {
	batchTimeout, err := time.ParseDuration(cfg.Kafka.BatchTimeout)
	if err != nil {
		return nil, err
	}
	dialerTimeout, err := time.ParseDuration(cfg.Kafka.DialerTimeout)
	if err != nil {
		return nil, err
	}

	return &producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{cfg.Kafka.Address},
			Topic:        topic,
			BatchTimeout: batchTimeout,
			//Balancer:     &kafka.RoundRobin{},
			Dialer: &kafka.Dialer{
				Timeout:   dialerTimeout,
				DualStack: true,
			},
		}),
	}, nil
}

func (m *producer) WriteMessages(ctx context.Context, value []byte) error {
	err := m.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(getUUID()),
			Value: value,
			Time:  time.Now(),
		},
	)
	if err != nil {
		return err
	}

	return m.writer.Close()
}

func (m *producer) Close() error {
	return m.writer.Close()
}

func getUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}