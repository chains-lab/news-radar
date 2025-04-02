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
	//UpdateName(ctx context.Context, name string, newName string) error

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
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tags) Create(input TagCreateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := t.neo.Create(ctxSync, neodb.TagCreateInput{
		Name:   input.Name,
		Status: input.Status,
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

func (t *Tags) Update(name string, fields map[string]any) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if color, ok := fields["color"].(string); ok {
		err := t.redis.UpdateColor(ctxSync, name, color)
		if err != nil {
			return err
		}
	}

	if icon, ok := fields["icon"].(string); ok {
		err := t.redis.UpdateIcon(ctxSync, name, icon)
		if err != nil {
			return err
		}
	}

	if status, ok := fields["status"].(string); ok {
		err := t.neo.UpdateStatus(ctxSync, name, status)
		if err != nil {
			return err
		}
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
		Status: neo.Status,
		Name:   redis.Name,
		Color:  redis.Color,
		Icon:   redis.Icon,
	}, nil
}
