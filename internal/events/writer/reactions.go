package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/events"
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

//func (p *Writer) Like(userID uuid.UUID, articleID uuid.UUID) error {
//	key := articleID.String()
//	payload := map[string]string{
//		"user_id":    userID.String(),
//		"article_id": articleID.String(),
//		"action":     strings.ToLower(events.LikeEventType),
//	}
//	body, err := json.Marshal(payload)
//	if err != nil {
//		return fmt.Errorf("failed to marshal like payload: %w", err)
//	}
//	return p.SendMessage("reactions", events.LikeEventType, key, body)
//}
//
//func (p *Writer) LikeRemove(userID uuid.UUID, articleID uuid.UUID) error {
//	key := articleID.String()
//	payload := map[string]string{
//		"user_id":    userID.String(),
//		"article_id": articleID.String(),
//		"action":     strings.ToLower(events.LikeRemoveEventType),
//	}
//	body, err := json.Marshal(payload)
//	if err != nil {
//		return fmt.Errorf("failed to marshal like_remove payload: %w", err)
//	}
//	return p.SendMessage("reactions", events.LikeRemoveEventType, key, body)
//}
//
//
//func (p *Writer) Repost(userID uuid.UUID, articleID uuid.UUID) error {
//	key := articleID.String()
//	payload := map[string]string{
//		"user_id":    userID.String(),
//		"article_id": articleID.String(),
//		"action":     strings.ToLower(events.RepostEventType),
//	}
//	body, err := json.Marshal(payload)
//	if err != nil {
//		return fmt.Errorf("failed to marshal repost payload: %w", err)
//	}
//	return p.SendMessage("reactions", events.RepostEventType, key, body)
//}
