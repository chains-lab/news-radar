package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type Tags interface {
}

type tags struct {
	data repo.Tags
}

func NewTags(cfg config.Config) (Tags, error) {
	repo, err := repo.NewTags(cfg)
	if err != nil {
		return nil, err
	}

	return &tags{
		data: repo,
	}, nil
}
