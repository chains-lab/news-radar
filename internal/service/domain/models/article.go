package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID     `json:"id" bson:"_id"`
	Title     string        `json:"title" bson:"title"`
	Icon      string        `json:"icon" bson:"icon"`
	Desc      string        `json:"desc" bson:"desc"`
	Authors   []uuid.UUID   `json:"authors,omitempty" bson:"authors,omitempty"`
	Content   []Section     `json:"content,omitempty" bson:"content,omitempty"`
	Likes     int           `json:"likes" bson:"likes"`
	Reposts   int           `json:"reposts" bson:"reposts"`
	Status    ArticleStatus `json:"status" bson:"status"`
	UpdatedAt *time.Time    `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

type Section struct {
	Section SectionType            `json:"section" bson:"section"`
	Content map[string]interface{} `json:"content" bson:"content"`
}

type SectionType string

const (
	SectionTypeText     SectionType = "text"
	SectionTypeImage    SectionType = "image"
	SectionTypeVideo    SectionType = "video"
	SectionTypeLocation SectionType = "location"
	SectionTypeQuote    SectionType = "quote"
)

func SectionTypeParse(s string) (SectionType, error) {
	switch s {
	case "text":
		return SectionTypeText, nil
	case "image":
		return SectionTypeImage, nil
	case "video":
		return SectionTypeVideo, nil
	case "location":
		return SectionTypeLocation, nil
	case "quote":
		return SectionTypeQuote, nil
	default:
		return "", fmt.Errorf("invalid section type %s", s)
	}
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
