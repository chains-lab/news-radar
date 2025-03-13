package listener

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/service"
	"github.com/recovery-flow/news-radar/internal/service/infra/events"
	"github.com/segmentio/kafka-go"
)

type TopicConfig struct {
	Topic      string
	ReplyTopic string
	Callback   func(ctx context.Context, svc *service.Service, m kafka.Message, evt events.InternalEvent) error
	OnSuccess  func(ctx context.Context, svc *service.Service, m kafka.Message, ie events.InternalEvent) error
	OnError    func(ctx context.Context, svc *service.Service, m kafka.Message, ie events.InternalEvent, err error)
}

var TopicsConfig = []TopicConfig{}
