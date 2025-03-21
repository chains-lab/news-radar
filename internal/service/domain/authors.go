package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/infra"
	"github.com/recovery-flow/news-radar/internal/service/infra/data/neo"
	"github.com/recovery-flow/news-radar/internal/service/models"
	"github.com/sirupsen/logrus"
)

type Authors interface {
	Create(ctx context.Context, name string) (*models.Author, error)
	Update(ctx context.Context, ID uuid.UUID, fields map[string]any) (*models.Author, error)
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Author, error)
}

type authors struct {
	Infra *infra.Infra
	log   *logrus.Logger
}

func (a *authors) Create(ctx context.Context, name string) (*models.Author, error) {
	authorID := uuid.New()

	if err := a.Infra.Neo.Authors.Create(ctx, &neo.Author{
		ID:   authorID,
		Name: name,
	}); err != nil {
		return nil, err
	}

	res, err := a.Infra.Mongo.Authors.Insert(ctx, &models.Author{
		ID:        authorID,
		Name:      name,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *authors) Update(ctx context.Context, ID uuid.UUID, fields map[string]any) (*models.Author, error) {
	if _, ok := fields["name"]; !ok {
		if err := a.Infra.Neo.Authors.Update(ctx, ID, fields["name"].(string)); err != nil {
			return nil, err
		}
	}

	res, err := a.Infra.Mongo.Authors.New().FiltersID(ID).Update(ctx, fields)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *authors) Delete(ctx context.Context, ID uuid.UUID) error {
	if err := a.Infra.Neo.Authors.Delete(ctx, ID); err != nil {
		return err
	}

	if err := a.Infra.Mongo.Authors.New().FiltersID(ID).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (a *authors) GetByID(ctx context.Context, ID uuid.UUID) (*models.Author, error) {
	res, err := a.Infra.Mongo.Authors.New().FiltersID(ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
