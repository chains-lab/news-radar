package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ThemeModels struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

const themesNamespace = "ThemesImpl"

type ThemesImpl struct {
	client *redis.Client
}

func NewThemes(addr, password string, DB int) *ThemesImpl {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       DB,
	})

	return &ThemesImpl{
		client: redisClient,
	}
}

func (t *ThemesImpl) Add(ctx context.Context, theme ThemeModels) error {
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

func (t *ThemesImpl) Get(ctx context.Context, theme string) (*ThemeModels, error) {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)

	data, err := t.client.HGetAll(ctx, nameKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting theme from Redis: %w", err)
	}

	return &ThemeModels{
		Name:  data["name"],
		Color: data["color"],
		Icon:  data["icon"],
	}, nil
}

func (t *ThemesImpl) Delete(ctx context.Context, theme string) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)

	if err := t.client.Del(ctx, nameKey).Err(); err != nil {
		return fmt.Errorf("error deleting theme from Redis: %w", err)
	}

	return nil
}

func (t *ThemesImpl) Drop(ctx context.Context) error {
	keys, err := t.client.Keys(ctx, fmt.Sprintf("%s:*", themesNamespace)).Result()
	if err != nil {
		return fmt.Errorf("error listing ThemesImpl keys: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	if err := t.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("error deleting ThemesImpl keys: %w", err)
	}
	return nil
}

func (t *ThemesImpl) UpdateIcon(ctx context.Context, theme string, icon string) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)
	if err := t.client.HSet(ctx, nameKey, "icon", icon).Err(); err != nil {
		return fmt.Errorf("error updating theme icon in Redis: %w", err)
	}
	return nil
}

func (t *ThemesImpl) UpdateColor(ctx context.Context, theme string, color string) error {
	nameKey := fmt.Sprintf("%s:name:%s", themesNamespace, theme)
	if err := t.client.HSet(ctx, nameKey, "color", color).Err(); err != nil {
		return fmt.Errorf("error updating theme color in Redis: %w", err)
	}
	return nil
}
