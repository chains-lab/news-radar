package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID     uuid.UUID    `json:"id" bson:"id"`
	Name   string       `json:"name" bson:"name"`
	Status AuthorStatus `json:"status" bson:"status"`

	Desc   *string `json:"desc" bson:"desc"`
	Avatar *string `json:"avatar,omitempty" bson:"avatar,omitempty"`

	Email    *string `json:"email,omitempty" bson:"email,omitempty"`
	Telegram *string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter  *string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

type AuthorShort struct {
	ID   uuid.UUID `json:"id" bson:"id"`
	Name string    `json:"name" bson:"name"`
	Icon string    `json:"icon" bson:"icon"`
}

type AuthorStatus string

const (
	AuthorStatusActive   AuthorStatus = "active"
	AuthorStatusInactive AuthorStatus = "inactive"
)

func ParseAuthorStatus(s string) (AuthorStatus, error) {
	switch s {
	case "active":
		return AuthorStatusActive, nil
	case "inactive":
		return AuthorStatusInactive, nil
	default:
		return "", fmt.Errorf("unknown author status: %s", s)
	}
}
