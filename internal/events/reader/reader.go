package reader

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/recovery-flow/news-radar/internal/events"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Reader struct {
	log    *logrus.Entry
	reader *kafka.Reader
}

func NewReader(log *logrus.Entry, r *kafka.Reader) *Reader {
	return &Reader{
		log:    log,
		reader: r,
	}
}

func (r *Reader) ListenChan(ctx context.Context) <-chan events.InternalEvent {
	out := make(chan events.InternalEvent)

	go func() {
		defer r.reader.Close()
		defer close(out)

		for {
			m, err := r.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(ctx.Err(), context.Canceled) {
					r.log.Info("Context canceled, stopping listener")
					return
				}
				r.log.WithError(err).Error("Error reading message")
				continue
			}

			var ie events.InternalEvent
			if err := json.Unmarshal(m.Value, &ie); err != nil {
				r.log.WithError(err).Error("Error unmarshalling InternalEvent")
				continue
			}

			select {
			case out <- ie:
			case <-ctx.Done():
				r.log.Info("Context canceled while sending message to channel")
				return
			}
		}
	}()

	return out
}
