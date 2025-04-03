package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
	"github.com/hs-zavet/news-radar/internal/repo/redisdb"
)

type TagModel struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

type tagsRedis interface {
	Create(ctx context.Context, input redisdb.TagCreateInput) error
	Delete(ctx context.Context, tag string) error

	Get(ctx context.Context, tag string) (redisdb.TagModel, error)

	UpdateIcon(ctx context.Context, tag string, icon string) error
	UpdateColor(ctx context.Context, tag string, color string) error

	Drop(ctx context.Context) error
}

type tagsNeo interface {
	Create(ctx context.Context, input neodb.TagCreateInput) error
	Delete(ctx context.Context, name string) error

	UpdateStatus(ctx context.Context, name string, status string) error
	UpdateName(ctx context.Context, name string, newName string) error
	UpdateType(ctx context.Context, name string, newType string) error

	Get(ctx context.Context, name string) (neodb.TagModel, error)
	Select(ctx context.Context) ([]neodb.TagModel, error)
}

type Tags struct {
	redis tagsRedis
	neo   tagsNeo
}

func NewTags(cfg config.Config) (*Tags, error) {
	neo, err := neodb.NewTags(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	redis := redisdb.NewTags(cfg.Database.Redis.Addr, cfg.Database.Redis.Password, cfg.Database.Redis.DB)
	return &Tags{
		neo:   neo,
		redis: redis,
	}, nil
}

type TagCreateInput struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tags) Create(input TagCreateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := t.neo.Create(ctxSync, neodb.TagCreateInput{
		Name:      input.Name,
		Status:    input.Status,
		Type:      input.Type,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	err = t.redis.Create(ctxSync, redisdb.TagCreateInput{
		Name:  input.Name,
		Color: input.Color,
		Icon:  input.Icon,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) Delete(name string) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := t.neo.Delete(ctxSync, name)
	if err != nil {
		return err
	}

	err = t.redis.Delete(ctxSync, name)
	if err != nil {
		return err
	}

	return nil
}

type TagUpdateInput struct {
	Color  *string `json:"color"`
	Icon   *string `json:"icon"`
	Status *string `json:"status"`
	Type   *string `json:"type"`
	Name   *string `json:"name"`
}

func (t *Tags) Update(name string, input TagUpdateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if input.Color != nil {
		err := t.redis.UpdateColor(ctxSync, name, *input.Color)
		if err != nil {
			return err
		}
	}

	if input.Icon != nil {
		err := t.redis.UpdateIcon(ctxSync, name, *input.Icon)
		if err != nil {
			return err
		}
	}

	if input.Status != nil {
		err := t.neo.UpdateStatus(ctxSync, name, *input.Status)
		if err != nil {
			return err
		}
	}

	if input.Type != nil {
		err := t.neo.UpdateType(ctxSync, name, *input.Type)
		if err != nil {
			return err
		}
	}

	if input.Name != nil {
		err := t.neo.UpdateName(ctxSync, name, *input.Name)
		if err != nil {
			return err
		}
		err = t.redis.Delete(ctxSync, name)
		if err != nil {
			return err
		}
		err = t.redis.Create(ctxSync, redisdb.TagCreateInput{
			Name:  *input.Name,
			Color: *input.Color,
			Icon:  *input.Icon,
		})
	}

	return nil
}

func (t *Tags) Get(name string) (TagModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	neoTag, err := t.neo.Get(ctxSync, name)
	if err != nil {
		return TagModel{}, err
	}

	redisTag, err := t.redis.Get(ctxSync, name)
	if err != nil {
		return TagModel{}, err
	}

	tag, err := TagsCreateModel(redisTag, neoTag)
	if err != nil {
		return TagModel{}, err
	}

	return tag, nil
}

func TagsCreateModel(redis redisdb.TagModel, neo neodb.TagModel) (TagModel, error) {
	if redis.Name != neo.Name {
		return TagModel{}, fmt.Errorf("redis and neo names do not match")
	}

	return TagModel{
		Name:      neo.Name,
		Status:    neo.Status,
		Type:      neo.Type,
		Color:     redis.Color,
		Icon:      redis.Icon,
		CreatedAt: neo.CreatedAt,
	}, nil
}
