package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/data/mongodb"
	"github.com/hs-zavet/news-radar/internal/data/neodb"
)

type Article struct {
	ID        uuid.UUID  `json:"id" bson:"_id"`
	Title     string     `json:"title" bson:"title"`
	Icon      string     `json:"icon" bson:"icon"`
	Desc      string     `json:"desc" bson:"desc"`
	Content   []Section  `json:"content,omitempty" bson:"content,omitempty"`
	Likes     int        `json:"likes" bson:"likes"`
	Reposts   int        `json:"reposts" bson:"reposts"`
	Status    string     `json:"status" bson:"status"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
}

type Section struct {
	Section string         `json:"section" bson:"section"`
	Content map[string]any `json:"content" bson:"content"`
}

func CreateArticleModel(mongo mongodb.ArticleModel, neo neodb.ArticleModel) (Article, error) {
	sections := make([]Section, len(mongo.Content))
	for _, sec := range mongo.Content {
		section := Section{
			Section: sec.Section,
			Content: sec.Content,
		}
		sections = append(sections, section)
	}

	if mongo.ID != neo.ID {
		return Article{}, fmt.Errorf("mongo and neo IDs do not match")
	}

	return Article{
		ID:        mongo.ID,
		Title:     mongo.Title,
		Icon:      mongo.Icon,
		Desc:      mongo.Desc,
		Content:   sections,
		Likes:     mongo.Likes,
		Reposts:   mongo.Reposts,
		Status:    neo.Status,
		UpdatedAt: mongo.UpdatedAt,
		CreatedAt: mongo.CreatedAt,
	}, nil
}
