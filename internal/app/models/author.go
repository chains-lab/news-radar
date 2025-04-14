package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/enums"
)

type Author struct {
	ID        uuid.UUID          `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Status    enums.AuthorStatus `json:"status" bson:"status"`
	Desc      *string            `json:"desc" bson:"desc"`
	Avatar    *string            `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email     *string            `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string            `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string            `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
