package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/content"
)

type Article struct {
	ID     uuid.UUID     `json:"id" bson:"_id"`
	Status ArticleStatus `json:"status" bson:"status"`

	Title   string            `json:"title" bson:"title"`
	Icon    *string           `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc    *string           `json:"desc,omitempty" bson:"desc,omitempty"`
	Content []content.Section `json:"content,omitempty" bson:"content,omitempty"`

	Likes   int `json:"likes" bson:"likes"`
	Reposts int `json:"reposts" bson:"reposts"`

	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`

	Authors []uuid.UUID `json:"authors,omitempty" bson:"authors,omitempty"`
	Tags    []string    `json:"tags,omitempty" bson:"tags,omitempty"`
}

type ArticleStatus string

const (
	ArticleStatusPending  ArticleStatus = "pending"
	ArticleStatusActive   ArticleStatus = "active"
	ArticleStatusInactive ArticleStatus = "inactive"
)

func ParseArticleStatus(s string) (ArticleStatus, error) {
	switch s {
	case "pending":
		return ArticleStatusPending, nil
	case "active":
		return ArticleStatusActive, nil
	case "inactive":
		return ArticleStatusInactive, nil
	default:
		return "", fmt.Errorf("invalid status %s", s)
	}
}
