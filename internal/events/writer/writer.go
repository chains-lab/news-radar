package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/chains-lab/news-radar/internal/config"
	"github.com/chains-lab/news-radar/internal/events"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	brokers net.Addr
	writer  *kafka.Writer
}

func NewWriter(cfg config.Config, RequiredAcks kafka.RequiredAcks) *Writer {
	return &Writer{
		brokers: kafka.TCP(cfg.Kafka.Brokers...),
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    1,
			BatchTimeout: 0,
			Async:        false,
			RequiredAcks: RequiredAcks,
		},
	}
}

func (p *Writer) SendMessage(topic string, event string, key string, body []byte) error {
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
