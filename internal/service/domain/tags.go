package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/domain/models"
	"github.com/recovery-flow/news-radar/internal/service/infra"
	"github.com/sirupsen/logrus"
)

type Tags interface {
	Create(ctx context.Context, name string) (*models.Tag, error)
	Update(ctx context.Context, ID string, fields map[string]interface{}) (*models.Tag, error)
	Delete(ctx context.Context, ID string) error

	GetByID(ctx context.Context, ID string) (*models.Tag, error)
}

type tags struct {
	Infra *infra.Infra
	log   logrus.Logger
}

func (t *tags) Create(ctx context.Context, name string) (*models.Tag, error) {
	tagID := uuid.New()

	if err := t.Infra.Neo.Tags.Create(ctx, &neo.Tag{
		ID:   tagID,
		Name: name,
	}); err != nil {
		return nil, err
	}

	res, err := t.Infra.Mongo.Tags.Insert(ctx, &models.Tag{
		ID:        tagID,
		Name:      name,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
