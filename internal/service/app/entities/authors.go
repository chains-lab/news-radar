package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/data"
)

type Authors interface {
}

type authors struct {
	data data.Authors
}

func NewAuthors(cfg config.Config) (Authors, error) {
	repo, err := data.NewAuthors(cfg)
	if err != nil {
		return nil, err
	}

	return &authors{
		data: repo,
	}, nil
}
