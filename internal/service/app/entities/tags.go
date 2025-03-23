package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/data"
)

type Tags interface {
}

type tags struct {
	data data.Tags
}

func NewTags(cfg config.Config) (Tags, error) {
	repo, err := data.NewTags(cfg)
	if err != nil {
		return nil, err
	}

	return &tags{
		data: repo,
	}, nil
}
