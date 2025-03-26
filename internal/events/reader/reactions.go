package reader

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/events"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Reactions interface {
	Like(eve events.Reaction, articleID uuid.UUID) error
	LikeRemove(eve events.Reaction, articleID uuid.UUID) error
	Dislike(eve events.Reaction, articleID uuid.UUID) error
	DislikeRemove(eve events.Reaction, articleID uuid.UUID) error
	Repost(eve events.Reaction, articleID uuid.UUID) error

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

func (r *reactions) Like(eve events.Reaction, articleID uuid.UUID) error {

	//res, err := r.app.UserLike(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s Like article %s", eve.UserID, articleID)
	return nil
}

func (r *reactions) LikeRemove(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserLikeRemove(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s Like article %s", eve.UserID, articleID)
	return nil
}

func (r *reactions) Dislike(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserDislike(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s Dislike article %s", eve.UserID, articleID)

	return nil
}

func (r *reactions) DislikeRemove(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserDislikeRemove(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s Dislike article %s", eve.UserID, articleID)
	return nil
}

func (r *reactions) Repost(eve events.Reaction, articleID uuid.UUID) error {

	//_, err = r.app.UserRepost(ctx, accountID, articleID)
	//if err != nil {
	//	return err
	//}

	r.log.Infof("user %s Repost article %s", eve.UserID, articleID)
	return nil
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
				if err := r.Like(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing Like reaction")
					continue
				}
			case events.LikeRemoveEventType:
				if err := r.LikeRemove(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing Like removal reaction")
					continue
				}
			case events.DislikeEventType:
				if err := r.Dislike(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing Dislike reaction")
					continue
				}
			case events.DislikeRemoveEventType:
				if err := r.DislikeRemove(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing Dislike removal reaction")
					continue
				}
			case events.RepostEventType:
				if err := r.Repost(eve, articleID); err != nil {
					r.log.WithError(err).Error("Error processing Repost reaction")
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
