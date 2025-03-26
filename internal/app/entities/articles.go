package entities

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/events/writer"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type Articles interface {
}

type articles struct {
	data   repo.Article
	events writer.Reaction
}

func NewArticles(cfg config.Config) (Articles, error) {
	repo, err := repo.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	events := writer.NewReactions(cfg)

	return &articles{
		data:   repo,
		events: events,
	}, nil
}
