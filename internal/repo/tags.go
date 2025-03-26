package repo

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
	"github.com/recovery-flow/news-radar/internal/data/redisdb"
)

type Tags interface {
	Create(ctx context.Context, tag models.Tag) error
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, name string, fields map[string]any) error
	Get(ctx context.Context, name string) (*models.Tag, error)
}

type tags struct {
	redis redisdb.Tags
	neo   neodb.Tags
}

func NewTags(cfg config.Config) (Tags, error) {
	neo, err := neodb.NewTags(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	redis := redisdb.NewTags(cfg.Database.Redis.Addr, cfg.Database.Redis.Password, cfg.Database.Redis.DB)
	return &tags{
		neo:   neo,
		redis: redis,
	}, nil
}

func (t *tags) Create(ctx context.Context, tag models.Tag) error {
	neoTag := neodb.Tag{
		Name:      tag.Name,
		Status:    tag.Status,
		CreatedAt: tag.CreatedAt,
	}
	err := t.neo.Create(ctx, neoTag)
	if err != nil {
		return err
	}

	err = t.redis.Add(ctx, redisdb.Tag{
		Name:  tag.Name,
		Color: tag.Color,
		Icon:  tag.Icon,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *tags) Delete(ctx context.Context, name string) error {
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

func (t *tags) Update(ctx context.Context, name string, fields map[string]any) error {
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

func (t *tags) Get(ctx context.Context, name string) (*models.Tag, error) {
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

func createModelsTag(neo neodb.Tag, redis redisdb.Tag) models.Tag {
	return models.Tag{
		Name:      neo.Name,
		Color:     redis.Color,
		Icon:      redis.Icon,
		Status:    neo.Status,
		CreatedAt: neo.CreatedAt,
	}
}
