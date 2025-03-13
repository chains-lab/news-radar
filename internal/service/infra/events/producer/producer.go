package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/infra/events"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
}

type producer struct {
	brokers net.Addr
	writer  *kafka.Writer
}

func NewProducer(cfg *config.Config) Producer {
	return &producer{
		brokers: kafka.TCP(cfg.Kafka.Brokers...),
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    1,
			BatchTimeout: 0,
			Async:        false,
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *producer) sendMessage(topic string, event string, key string, body []byte) error {
	evt := events.InternalEvent{
		EventType: event,
		Data:      body,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription activate event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Value: data,
		Key:   []byte(key),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}
