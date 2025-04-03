package models

import (
	"fmt"
	"time"
)

type Tag struct {
	Name      string    `json:"name"`
	Status    TagStatus `json:"status"`
	Type      TagType   `json:"type"`
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

type TagType string

const (
	TagTypeTopic    TagType = "topic"
	TagTypeSubTopic TagType = "sub-topic"
	TagTypeDefault  TagType = "default"
)

func ParseTagType(s string) (TagType, error) {
	switch s {
	case "topic":
		return TagTypeTopic, nil
	case "sub-topic":
		return TagTypeSubTopic, nil
	case "default":
		return TagTypeDefault, nil
	default:
		return "", fmt.Errorf("invalid tag type %s", s)
	}
}
