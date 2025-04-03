package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	Type      ReactionType `json:"type"`
	UserID    uuid.UUID    `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
}

type ReactionType string

const (
	ReactionTypeLike   ReactionType = "like"
	ReactionTypeRepost ReactionType = "repost"
)

func ParseReactionType(s string) (ReactionType, error) {
	switch s {
	case "like":
		return ReactionTypeLike, nil
	case "repost":
		return ReactionTypeRepost, nil
	default:
		return "", fmt.Errorf("invalid reaction type: %s", s)
	}
}
