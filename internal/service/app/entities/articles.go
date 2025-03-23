package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/data"
)

type Articles interface {
}

type articles struct {
	data data.Article
}

func NewArticles(cfg config.Config) (Articles, error) {
	repo, err := data.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	return &articles{
		data: repo,
	}, nil
}
