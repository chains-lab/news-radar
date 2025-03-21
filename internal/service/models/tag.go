package models

import (
	"fmt"
	"time"
)

type Tag struct {
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

func ParseTagType(s string) (TagType, error) {
	switch s {
	case "category":
		return TagTypeCategory, nil
	case "person":
		return TagTypeTag, nil
	default:
		return "", fmt.Errorf("invalid tag type %s", s)
	}
}

type TagStatus string

const (
	TagStatusActive   TagStatus = "active"
	TagStatusInactive TagStatus = "inactive"
)

func ParseTagStatus(s string) (TagStatus, error) {
	switch s {
	case "active":
		return TagStatusActive, nil
	case "inactive":
		return TagStatusInactive, nil
	default:
		return "", fmt.Errorf("invalid tag status %s", s)
	}
}
