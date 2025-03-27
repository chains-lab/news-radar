package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data"
)

type authorsRepo interface {
	Create(ctx context.Context, author models.Author) error
	Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Author, error)
}

type Authors struct {
	data authorsRepo
}

func NewAuthors(cfg config.Config) (*Authors, error) {
	repo, err := data.NewAuthors(cfg)
	if err != nil {
		return nil, err
	}

	return &Authors{
		data: repo,
	}, nil
}
