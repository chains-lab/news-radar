package models

import (
	"time"

	"github.com/chains-lab/news-radar/internal/content"
	"github.com/chains-lab/news-radar/internal/enums"
	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID           `json:"id" bson:"_id"`
	Status      enums.ArticleStatus `json:"status" bson:"status"`
	Title       string              `json:"title" bson:"title"`
	Icon        *string             `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc        *string             `json:"desc,omitempty" bson:"desc,omitempty"`
	Content     []content.Section   `json:"content,omitempty" bson:"content,omitempty"`
	PublishedAt *time.Time          `json:"published_at,omitempty" bson:"published_at,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
}
