package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/content"
)

type Article struct {
	ID        uuid.UUID         `json:"id" bson:"_id"`
	Title     string            `json:"title" bson:"title"`
	Icon      string            `json:"icon" bson:"icon"`
	Desc      string            `json:"desc" bson:"desc"`
	Content   []content.Section `json:"content,omitempty" bson:"content,omitempty"`
	Likes     int               `json:"likes" bson:"likes"`
	Reposts   int               `json:"reposts" bson:"reposts"`
	Status    ArticleStatus     `json:"status" bson:"status"`
	UpdatedAt *time.Time        `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
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
