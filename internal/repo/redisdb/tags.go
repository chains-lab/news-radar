package redisdb

import (
	"context"
	"fmt"

	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
	"github.com/redis/go-redis/v9"
)

const tagsNamespace = "tags"

type TagsImpl struct {
	client *redis.Client
}

func NewTags(addr, password string, DB int) *TagsImpl {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})

	return &TagsImpl{
		client: redisClient,
	}
}

func (t *TagsImpl) Add(ctx context.Context, tag modelsdb.TagRedis) error {
	nameKey := fmt.Sprintf("%s:name:%s", tagsNamespace, tag.Name)

	data := map[string]interface{}{
		"name":  tag.Name,
		"color": tag.Color,
		"icon":  tag.Icon,
	}

	if err := t.client.HSet(ctx, nameKey, data).Err(); err != nil {
		return fmt.Errorf("error adding tag to Redis: %w", err)
	}

	return nil
}

func (t *TagsImpl) Get(ctx context.Context, tag string) (modelsdb.TagRedis, error) {
	nameKey := fmt.Sprintf("%s:name:%s", tagsNamespace, tag)

	data, err := t.client.HGetAll(ctx, nameKey).Result()
	if err != nil {
		return modelsdb.TagRedis{}, fmt.Errorf("error getting tag from Redis: %w", err)
	}

	return modelsdb.TagRedis{
		Name:  data["name"],
		Color: data["color"],
		Icon:  data["icon"],
	}, nil
}

func (t *TagsImpl) Delete(ctx context.Context, tag string) error {
	nameKey := fmt.Sprintf("%s:name:%s", tagsNamespace, tag)

	if err := t.client.Del(ctx, nameKey).Err(); err != nil {
		return fmt.Errorf("error deleting tag from Redis: %w", err)
	}

	return nil
}

func (t *TagsImpl) Drop(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", tagsNamespace)

	keys, err := t.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("error getting keys from Redis: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	if err := t.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("error deleting tags from Redis: %w", err)
	}

	return nil
}

func (t *TagsImpl) UpdateIcon(ctx context.Context, tag string, icon string) error {
	nameKey := fmt.Sprintf("%s:name:%s", tagsNamespace, tag)

	if err := t.client.HSet(ctx, nameKey, "icon", icon).Err(); err != nil {
		return fmt.Errorf("error updating tag icon in Redis: %w", err)
	}

	return nil
}

func (t *TagsImpl) UpdateColor(ctx context.Context, tag string, color string) error {
	nameKey := fmt.Sprintf("%s:name:%s", tagsNamespace, tag)

	if err := t.client.HSet(ctx, nameKey, "color", color).Err(); err != nil {
		return fmt.Errorf("error updating tag color in Redis: %w", err)
	}

	return nil
}
