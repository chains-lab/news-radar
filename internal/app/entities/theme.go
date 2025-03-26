package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type Themes interface {
}

type theme struct {
	data repo.Themes
}

func NewThemes(cfg config.Config) (Themes, error) {
	repo, err := repo.NewThemes(cfg)
	if err != nil {
		return nil, err
	}

	return &theme{
		data: repo,
	}, nil
}
