package models

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID   `json:"id" bson:"_id"`
	Title       string      `json:"title" bson:"title"`
	Icon        string      `json:"icon" bson:"icon"`
	Description string      `json:"description" bson:"description"`
	Authors     []uuid.UUID `json:"authors" bson:"authors"`
	Content     Section     `json:"content" bson:"content"`
	Likes       int         `json:"likes" bson:"likes"`
	Reposts     int         `json:"reposts" bson:"reposts"`
	Tags        []Tag       `json:"tags" bson:"tags"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at"`
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
