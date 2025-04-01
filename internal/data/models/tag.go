package models

import (
	"fmt"
	"time"

	"github.com/hs-zavet/news-radar/internal/data/neodb"
	"github.com/hs-zavet/news-radar/internal/data/redisdb"
)

type Tag struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
}

func TagsCreateModel(redis redisdb.TagModels, neo neodb.TagModels) (Tag, error) {
	if redis.Name != neo.Name {
		return Tag{}, fmt.Errorf("redis and neo names do not match")
	}

	return Tag{
		Status:    neo.Status,
		CreatedAt: neo.CreatedAt,
		Name:      redis.Name,
		Color:     redis.Color,
		Icon:      redis.Icon,
	}, nil
}
