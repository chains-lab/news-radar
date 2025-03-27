package entities

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data"
)

type themesRepo interface {
	Create(ctx context.Context, theme models.Theme) error
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, name string, fields map[string]any) error
	Get(ctx context.Context, name string) (*models.Theme, error)
}

type Theme struct {
	data themesRepo
}

func NewThemes(cfg config.Config) (*Theme, error) {
	repo, err := data.NewThemes(cfg)
	if err != nil {
		return nil, err
	}

	return &Theme{
		data: repo,
	}, nil
}
