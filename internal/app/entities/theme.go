package entities

import (
	"context"

	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type ThemesRepo interface {
	Create(ctx context.Context, theme models.Theme) error
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, name string, fields map[string]any) error
	Get(ctx context.Context, name string) (*models.Theme, error)
}

type Theme struct {
	data ThemesRepo
}

func NewThemes(cfg config.Config) (*Theme, error) {
	repo, err := repo.NewThemes(cfg)
	if err != nil {
		return nil, err
	}

	return &Theme{
		data: repo,
	}, nil
}
