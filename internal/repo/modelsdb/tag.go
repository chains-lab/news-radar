package modelsdb

import (
	"fmt"
	"time"
)

type Tag struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

// TagRedis is the model for the tag in the Redis
type TagRedis struct {
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

// TagNeo is the model for the tag in the Neo4j
type TagNeo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func TagsCreateModel(redis TagRedis, neo TagNeo) (Tag, error) {
	if redis.Name != neo.Name {
		return Tag{}, fmt.Errorf("redis and neo names do not match")
	}

	return Tag{
		Status: neo.Status,
		Name:   redis.Name,
		Color:  redis.Color,
		Icon:   redis.Icon,
	}, nil
}
