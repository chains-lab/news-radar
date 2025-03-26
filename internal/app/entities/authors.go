package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type Authors interface {
}

type authors struct {
	data repo.Authors
}

func NewAuthors(cfg config.Config) (Authors, error) {
	repo, err := repo.NewAuthors(cfg)
	if err != nil {
		return nil, err
	}

	return &authors{
		data: repo,
	}, nil
}
