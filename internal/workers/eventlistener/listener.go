package eventlistener

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/config"
)

func Listen(ctx context.Context, cfg *config.Config, app app.App) {
	logger := cfg.Log().WithField("listener", "kafka")

	//reactionsWriter := reader.NewReactions(logger, app, kafka.NewReader(kafka.ReaderConfig{
	//	Brokers:        cfg.Kafka.Brokers,
	//	Topic:          events.ReactionsTopic,
	//	MinBytes:       1,
	//	MaxBytes:       10e6,
	//	CommitInterval: time.Second,
	//}))
	//
	//accountsWriter := reader.NewReader(logger, app, kafka.NewReader(kafka.ReaderConfig{
	//	Brokers:        cfg.Kafka.Brokers,
	//	Topic:          events.AccountsTopic,
	//	MinBytes:       1,
	//	MaxBytes:       10e6,
	//	CommitInterval: time.Second,
	//}))
	//
	//go reactionsWriter.Listen(ctx)
	//go accountsWriter.Listen(ctx)

	<-ctx.Done()
	logger.Info("Producer listener stopped")
}
