package models

import (
	"fmt"
	"time"
)

type Tag struct {
	Name      string    `json:"name"`
	Status    TagStatus `json:"status"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
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
