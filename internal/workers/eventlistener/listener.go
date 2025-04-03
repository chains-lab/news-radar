package eventlistener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/events"
	"github.com/hs-zavet/news-radar/internal/events/reader"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Listener struct {
	cfg *config.Config
	app *app.App
	log *logrus.Entry
}

func NewListener(cfg *config.Config, app *app.App) *Listener {
	return &Listener{
		cfg: cfg,
		app: app,
		log: cfg.Log.WithField("listener", "kafka"),
	}
}

func (l *Listener) Listen(ctx context.Context, cfg *config.Config) {
	logger := cfg.Log.WithField("listener", "kafka")

	reactionReader := reader.NewReader(l.log, kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          events.ReactionsTopic,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	}))

	reactionChanel := reactionReader.ListenChan(ctx)

	go func(ctx context.Context) {
		for event := range reactionChanel {
			var eve events.Reaction
			if err := json.Unmarshal(event.Data, &eve); err != nil {
				l.log.WithError(err).Error("Error unmarshalling reaction event")
				continue
			}

			switch event.EventType {
			//case events.RepostEventType:
			//	l.app.Repost(ctx, event)
			//case events.LikeEventType:
			//	l.app.Like(ctx, event)
			//case events.LikeRemoveEventType:
			//	l.app.LikeRemove(ctx, event)
			default:
				l.log.WithField("event", event).Error("Unknown event type")
			}
		}
	}(ctx)

	accountReader := reader.NewReader(l.log, kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          events.AccountsTopic,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	}))

	accountChanel := accountReader.ListenChan(ctx)

	go func(ctx context.Context) {
		for event := range accountChanel {
			var eve events.AccountCreated
			if err := json.Unmarshal(event.Data, &eve); err != nil {
				l.log.WithError(err).Error("Error unmarshalling account create event")
				continue
			}

			switch event.EventType {
			//case events.AccountCreateType:
			//	l.app.AccountCreated(ctx, event)
			default:
				l.log.WithField("event", event).Error("Unknown event type")
			}
		}
	}(ctx)

	<-ctx.Done()
	logger.Info("Producer listener stopped")
}
