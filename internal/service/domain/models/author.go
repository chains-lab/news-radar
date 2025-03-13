package models

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID  `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Desc      string     `json:"desc" bson:"desc"`
	Avatar    string     `json:"avatar" bson:"avatar"`
	Email     *string    `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string    `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string    `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}
