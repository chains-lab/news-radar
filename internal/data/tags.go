package data

import (
	"context"

	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/data/models"
	"github.com/hs-zavet/news-radar/internal/data/neodb"
	"github.com/hs-zavet/news-radar/internal/data/redisdb"
)

type tagsRedis interface {
	Add(ctx context.Context, tag redisdb.TagModels) error
	Get(ctx context.Context, tag string) (redisdb.TagModels, error)
	Delete(ctx context.Context, tag string) error

	UpdateIcon(ctx context.Context, tag string, icon string) error
	UpdateColor(ctx context.Context, tag string, color string) error

	Drop(ctx context.Context) error
}

type tagsNeo interface {
	Create(ctx context.Context, tag neodb.TagModels) error
	Delete(ctx context.Context, name string) error

	UpdateStatus(ctx context.Context, name string, status string) error
	//UpdateName(ctx context.Context, name string, newName string) error

	Get(ctx context.Context, name string) (neodb.TagModels, error)
	Select(ctx context.Context) ([]neodb.TagModels, error)
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

func (t *Tags) Create(tag models.Tag) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	neoTag := neodb.TagModels{
		Name:   tag.Name,
		Status: tag.Status,
	}
	err := t.neo.Create(ctxSync, neoTag)
	if err != nil {
		return err
	}

	err = t.redis.Add(ctxSync, redisdb.TagModels{
		Name:  tag.Name,
		Color: tag.Color,
		Icon:  tag.Icon,
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

func (t *Tags) Get(name string) (*models.Tag, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	neoTag, err := t.neo.Get(ctxSync, name)
	if err != nil {
		return nil, err
	}

	redisTag, err := t.redis.Get(ctxSync, name)
	if err != nil {
		return nil, err
	}

	tag, err := models.TagsCreateModel(redisTag, neoTag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}
