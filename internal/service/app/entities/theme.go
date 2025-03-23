package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/data"
)

type Themes interface {
}

type theme struct {
	data data.Themes
}

func NewThemes(cfg config.Config) (Themes, error) {
	repo, err := data.NewThemes(cfg)
	if err != nil {
		return nil, err
	}

	return &theme{
		data: repo,
	}, nil
}
