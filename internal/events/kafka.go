package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	ReactionsTopic = "reactions"
	AccountsTopic  = "accounts"

	LikeEventType          = "LIKE"
	LikeRemoveEventType    = "LIKE_REMOVE"
	DislikeEventType       = "DISLIKE"
	DislikeRemoveEventType = "DISLIKE_REMOVE"
	RepostEventType        = "REPOST"
	AccountCreateType      = "ACCOUNT_CREATE"
)

type InternalEvent struct {
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

type Reaction struct {
	UserID    uuid.UUID `json:"user_id"`
	ArticleID uuid.UUID `json:"article_id"`
	Timestamp time.Time `json:"timestamp"`
}

type AccountCreated struct {
	AccountID string    `json:"account_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}
