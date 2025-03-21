package listener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/events"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Listener interface {
}

type listener struct {
	topics []TopicConfig
	log    *logrus.Logger
}

type TopicConfig struct {
	Topic      string
	ReplyTopic string
	Callback   func(ctx context.Context, m kafka.Message, e events.InternalEvent) error
	OnSuccess  func(ctx context.Context, m kafka.Message, e events.InternalEvent) error
	OnError    func(ctx context.Context, m kafka.Message, e events.InternalEvent, err error)
}

func NewListener(topics []TopicConfig, log *logrus.Logger) (Listener, error) {
	return &listener{
		topics: topics,
		log:    log,
	}, nil
}

func (l *listener) Listen(ctx context.Context, cfg config.Config) {
	for _, tc := range l.topics {
		tc := tc
		go func() {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:        cfg.Kafka.Brokers,
				Topic:          tc.Topic,
				MinBytes:       1,
				MaxBytes:       10e6,
				CommitInterval: time.Second,
			})
			defer r.Close()

			for {
				m, err := r.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					l.log.WithField("kafka", err).Errorf("Error reading message from topic %s", tc.Topic)
					continue
				}

				var ie events.InternalEvent
				if err := json.Unmarshal(m.Value, &ie); err != nil {
					l.log.WithField("kafka", err).Error("Error unmarshalling InternalEvent")
					continue
				}

				if err := tc.Callback(ctx, m, ie); err != nil {
					l.log.WithField("kafka", err).Errorf("Error processing message from topic %s", tc.Topic)
					if tc.OnError != nil {
						tc.OnError(ctx, m, ie, err)
					}
					continue
				}

				if tc.OnSuccess != nil {
					if err := tc.OnSuccess(ctx, m, ie); err != nil {
						l.log.WithField("kafka", err).Error("Error in OnSuccess callback")
					}
				}
			}
		}()
	}

	<-ctx.Done()
	l.log.Info("Producer listener stopped")
}
