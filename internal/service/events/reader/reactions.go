package reader

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/app"
	"github.com/recovery-flow/news-radar/internal/service/events"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Reactions interface {
	like(eve events.Reaction, articleID uuid.UUID) error
	likeRemove(eve events.Reaction, articleID uuid.UUID) error
	dislike(eve events.Reaction, articleID uuid.UUID) error
	dislikeRemove(eve events.Reaction, articleID uuid.UUID) error
	repost(eve events.Reaction, articleID uuid.UUID) error

	Listen(ctx context.Context)
}

type reactions struct {
	log *logrus.Entry
	app app.App
	*kafka.Reader
}

func NewReactions(log *logrus.Entry, app app.App, r *kafka.Reader) Reactions {
	return &reactions{
		log:    log,
		app:    app,
		Reader: r,
	}
}

func (r *reactions) Listen(ctx context.Context) {
	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					continue
				}
				r.log.WithError(err).Errorf("Error reading message from topic %s", events.ReactionsTopic)
				continue
			}

			var ie events.InternalEvent
			if err := json.Unmarshal(m.Value, &ie); err != nil {
				r.log.WithError(err).Error("Error unmarshalling InternalEvent")
				continue
			}

			var eve events.Reaction
			if err := json.Unmarshal(ie.Data, &eve); err != nil {
				r.log.WithError(err).Error("Error unmarshalling Reaction")
				continue
			}

			articleID, err := uuid.Parse(string(m.Key))
			if err != nil {
				r.log.WithError(err).Error("Error parsing article ID")
				continue
			}

			switch ie.EventType {
			case events.LikeEventType:
				if err := r.like(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing like reaction")
					continue
				}
			case events.LikeRemoveEventType:
				if err := r.likeRemove(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing like removal reaction")
					continue
				}
			case events.DislikeEventType:
				if err := r.dislike(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing dislike reaction")
					continue
				}
			case events.DislikeRemoveEventType:
				if err := r.dislikeRemove(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing dislike removal reaction")
					continue
				}
			case events.RepostEventType:
				if err := r.repost(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing repost reaction")
					continue
				}
			default:
				r.log.Errorf("Unknown event type %s", ie.EventType)
			}
		}
	}()
	<-ctx.Done()
	r.log.Info("Reactions listener stopped")
}

func (r *reactions) like(eve events.Reaction, articleID uuid.UUID) error {

	//res, err := r.app.UserLike(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s like article %s", eve.UserID, articleID)
	return nil
}

func (r *reactions) likeRemove(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserLikeRemove(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s like article %s", eve.UserID, articleID)
	return nil
}

func (r *reactions) dislike(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserDislike(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s dislike article %s", eve.UserID, articleID)

	return nil
}

func (r *reactions) dislikeRemove(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserDislikeRemove(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s dislike article %s", eve.UserID, articleID)
	return nil
}

func (r *reactions) repost(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserRepost(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s repost article %s", eve.UserID, articleID)
	return nil
}
