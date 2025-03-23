package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/events"

	"github.com/segmentio/kafka-go"
)

type Reaction interface {
	Like(userID uuid.UUID, articleID uuid.UUID) error
	LikeRemove(userID uuid.UUID, articleID uuid.UUID) error

	Dislike(userID uuid.UUID, articleID uuid.UUID) error
	DislikeRemove(userID uuid.UUID, articleID uuid.UUID) error

	Repost(userID uuid.UUID, articleID uuid.UUID) error
}

type reaction struct {
	brokers net.Addr
	writer  *kafka.Writer
}

func NewReactions(cfg config.Config) Reaction {
	return &reaction{
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

func (p *reaction) sendMessage(topic string, event string, key string, body []byte) error {
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

func (p *reaction) Like(userID uuid.UUID, articleID uuid.UUID) error {
	key := articleID.String()
	payload := map[string]string{
		"user_id":    userID.String(),
		"article_id": articleID.String(),
		"action":     strings.ToLower(events.LikeEventType),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal like payload: %w", err)
	}
	return p.sendMessage("reactions", events.LikeEventType, key, body)
}

func (p *reaction) LikeRemove(userID uuid.UUID, articleID uuid.UUID) error {
	key := articleID.String()
	payload := map[string]string{
		"user_id":    userID.String(),
		"article_id": articleID.String(),
		"action":     strings.ToLower(events.LikeRemoveEventType),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal like_remove payload: %w", err)
	}
	return p.sendMessage("reactions", events.LikeRemoveEventType, key, body)
}

func (p *reaction) Dislike(userID uuid.UUID, articleID uuid.UUID) error {
	key := articleID.String()
	payload := map[string]string{
		"user_id":    userID.String(),
		"article_id": articleID.String(),
		"action":     strings.ToLower(events.DislikeEventType),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal dislike payload: %w", err)
	}
	return p.sendMessage("reactions", events.DislikeEventType, key, body)
}

func (p *reaction) DislikeRemove(userID uuid.UUID, articleID uuid.UUID) error {
	key := articleID.String()
	payload := map[string]string{
		"user_id":    userID.String(),
		"article_id": articleID.String(),
		"action":     strings.ToLower(events.DislikeRemoveEventType),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal dislike_remove payload: %w", err)
	}
	return p.sendMessage("reactions", events.DislikeRemoveEventType, key, body)
}

func (p *reaction) Repost(userID uuid.UUID, articleID uuid.UUID) error {
	key := articleID.String()
	payload := map[string]string{
		"user_id":    userID.String(),
		"article_id": articleID.String(),
		"action":     strings.ToLower(events.RepostEventType),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal repost payload: %w", err)
	}
	return p.sendMessage("reactions", events.RepostEventType, key, body)
}
