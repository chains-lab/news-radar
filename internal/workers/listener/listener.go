package listener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/recovery-flow/news-radar/internal/service"
	"github.com/recovery-flow/news-radar/internal/service/infra/events"
	"github.com/segmentio/kafka-go"
)

func Listen(ctx context.Context, svc *service.Service) {
	for _, tc := range TopicsConfig {
		tc := tc
		go func() {
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers:        svc.Config.Kafka.Brokers,
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
					//svc.Log.WithField("kafka", err).Errorf("Error reading message from topic %s", tc.Topic)
					continue
				}

				var ie events.InternalEvent
				if err := json.Unmarshal(m.Value, &ie); err != nil {
					svc.Log.WithField("kafka", err).Error("Error unmarshalling InternalEvent")
					continue
				}

				if err := tc.Callback(ctx, svc, m, ie); err != nil {
					svc.Log.WithField("kafka", err).Errorf("Error processing message from topic %s", tc.Topic)
					if tc.OnError != nil {
						tc.OnError(ctx, svc, m, ie, err)
					}
					continue
				}

				if tc.OnSuccess != nil {
					if err := tc.OnSuccess(ctx, svc, m, ie); err != nil {
						svc.Log.WithField("kafka", err).Error("Error in OnSuccess callback")
					}
				}
			}
		}()
	}

	<-ctx.Done()
	svc.Log.Info("Producer listener stopped")
}
