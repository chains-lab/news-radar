package models

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Status    TagStatus `json:"status"`
	Type      TagType   `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type TagType string

const (
	TagTypeCategory TagType = "category"
	TagTypeTag      TagType = "person"
)

type TagStatus string

const (
	TagStatusActive   TagStatus = "active"
	TagStatusInactive TagStatus = "inactive"
)
