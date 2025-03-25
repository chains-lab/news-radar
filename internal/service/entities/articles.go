package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/repo"
)

type Articles interface {
}

type articles struct {
	data repo.Article
}

func NewArticles(cfg config.Config) (Articles, error) {
	data, err := repo.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	return &articles{
		data: data,
	}, nil
}
