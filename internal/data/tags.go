package data

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
	"github.com/recovery-flow/news-radar/internal/data/redisdb"
)

type tagsRedis interface {
	Add(ctx context.Context, tag redisdb.TagModels) error
	Get(ctx context.Context, tag string) (*redisdb.TagModels, error)
	Delete(ctx context.Context, tag string) error

	UpdateIcon(ctx context.Context, tag string, icon string) error
	UpdateColor(ctx context.Context, tag string, color string) error

	Drop(ctx context.Context) error
}

type tagsNeo interface {
	Create(ctx context.Context, tag neodb.TagModels) error
	Delete(ctx context.Context, name string) error

	UpdateStatus(ctx context.Context, name string, status models.TagStatus) error
	//UpdateName(ctx context.Context, name string, newName string) error

	Get(ctx context.Context, name string) (*neodb.TagModels, error)
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

func (t *Tags) Create(ctx context.Context, tag models.Tag) error {
	neoTag := neodb.TagModels{
		Name:      tag.Name,
		Status:    tag.Status,
		CreatedAt: tag.CreatedAt,
	}
	err := t.neo.Create(ctx, neoTag)
	if err != nil {
		return err
	}

	err = t.redis.Add(ctx, redisdb.TagModels{
		Name:  tag.Name,
		Color: tag.Color,
		Icon:  tag.Icon,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) Delete(ctx context.Context, name string) error {
	err := t.neo.Delete(ctx, name)
	if err != nil {
		return err
	}

	err = t.redis.Delete(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tags) Update(ctx context.Context, name string, fields map[string]any) error {
	if _, ok := fields["color"]; ok {
		err := t.redis.UpdateColor(ctx, name, fields["color"].(string))
		if err != nil {
			return err
		}
	}

	if _, ok := fields["icon"]; ok {
		err := t.redis.UpdateIcon(ctx, name, fields["icon"].(string))
		if err != nil {
			return err
		}
	}

	if _, ok := fields["status"]; ok {
		status, err := models.ParseTagStatus(fields["status"].(string))
		if err != nil {
			return err
		}
		err = t.neo.UpdateStatus(ctx, name, status)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Tags) Get(ctx context.Context, name string) (*models.Tag, error) {
	neoTag, err := t.neo.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	redisTag, err := t.redis.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	tag := createModelsTag(*neoTag, *redisTag)

	return &tag, nil
}

func createModelsTag(neo neodb.TagModels, redis redisdb.TagModels) models.Tag {
	return models.Tag{
		Name:      neo.Name,
		Color:     redis.Color,
		Icon:      redis.Icon,
		Status:    neo.Status,
		CreatedAt: neo.CreatedAt,
	}
}
