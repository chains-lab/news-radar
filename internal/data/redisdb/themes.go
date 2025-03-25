package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Theme struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type Themes interface {
	Add(ctx context.Context, theme Theme) error
	Get(ctx context.Context, theme string) (*Theme, error)
	Delete(ctx context.Context, theme string) error

	UpdateIcon(ctx context.Context, theme string, icon string) error
	UpdateColor(ctx context.Context, theme string, color string) error

	Drop(ctx context.Context) error
}

const themesNamespace = "themes"

type themes struct {
	client *redis.Client
}

func NewThemes(addr, password string, DB int) Themes {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})

	return &themes{
		client: redisClient,
	}
}

func (t *themes) Add(ctx context.Context, theme Theme) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme.Name)

	data := map[string]interface{}{
		"name":  theme.Name,
		"color": theme.Color,
		"icon":  theme.Icon,
	}

	if err := t.client.HSet(ctx, nameKey, data).Err(); err != nil {
		return fmt.Errorf("error adding theme to Redis: %w", err)
	}

	return nil
}

func (t *themes) Get(ctx context.Context, theme string) (*Theme, error) {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)

	data, err := t.client.HGetAll(ctx, nameKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting theme from Redis: %w", err)
	}

	return &Theme{
		Name:  data["name"],
		Color: data["color"],
		Icon:  data["icon"],
	}, nil
}

func (t *themes) Delete(ctx context.Context, theme string) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)

	if err := t.client.Del(ctx, nameKey).Err(); err != nil {
		return fmt.Errorf("error deleting theme from Redis: %w", err)
	}

	return nil
}

func (t *themes) Drop(ctx context.Context) error {
	keys, err := t.client.Keys(ctx, fmt.Sprintf("%s:*", themesNamespace)).Result()
	if err != nil {
		return fmt.Errorf("error listing themes keys: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	if err := t.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("error deleting themes keys: %w", err)
	}
	return nil
}

func (t *themes) UpdateIcon(ctx context.Context, theme string, icon string) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)
	if err := t.client.HSet(ctx, nameKey, "icon", icon).Err(); err != nil {
		return fmt.Errorf("error updating theme icon in Redis: %w", err)
	}
	return nil
}

func (t *themes) UpdateColor(ctx context.Context, theme string, color string) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)
	if err := t.client.HSet(ctx, nameKey, "color", color).Err(); err != nil {
		return fmt.Errorf("error updating theme color in Redis: %w", err)
	}
	return nil
}
