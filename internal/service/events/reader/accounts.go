package reader

import (
	"context"
	"encoding/json"

	"github.com/recovery-flow/news-radar/internal/service/app"
	"github.com/recovery-flow/news-radar/internal/service/events"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Accounts interface {
	Create(eve events.AccountCreated) error

	Listen(ctx context.Context)
}

type accounts struct {
	log *logrus.Entry
	app app.App
	*kafka.Reader
}

func NewAccounts(log *logrus.Entry, app app.App, r *kafka.Reader) Accounts {
	return &accounts{
		log:    log,
		app:    app,
		Reader: r,
	}
}

func (a *accounts) Create(eve events.AccountCreated) error {
	//res, err := a.app.Accounts.Create(ctx, eve.AccountID)
	//if err != nil {
	//	return err
	//}

	a.log.WithField("account", eve.AccountID).Info("Account created")
	return nil
}

func (a *accounts) Listen(ctx context.Context) {
	go func() {
		defer a.Close()
		for {
			m, err := a.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					continue
				}
				a.log.WithError(err).Errorf("Error reading message from topic %s", events.ReactionsTopic)
				continue
			}

			var ie events.InternalEvent
			if err := json.Unmarshal(m.Value, &ie); err != nil {
				a.log.WithError(err).Error("Error unmarshalling InternalEvent")
				continue
			}

			var eve events.AccountCreated
			if err := json.Unmarshal(ie.Data, &eve); err != nil {
				a.log.WithError(err).Error("Error unmarshalling AccountCreated")
				continue
			}

			switch ie.EventType {
			case events.AccountCreateType:
				if err := a.Create(eve); err != nil {
					a.log.WithError(err).Error("Error processing like reaction")
					continue
				}
			default:
				a.log.Errorf("Unknown event type %s", ie.EventType)
			}
		}
	}()
	<-ctx.Done()
	a.log.Info("Accounts listener stopped")
}
