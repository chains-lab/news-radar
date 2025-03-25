package service

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/entities"
)

type Domain interface {
}

type domain struct {
	articles entities.Articles
	authors  entities.Authors
	tags     entities.Tags
	themes   entities.Themes
	user     entities.User
}

func NewApp(cfg config.Config) (Domain, error) {
	articles, err := entities.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	authors, err := entities.NewAuthors(cfg)
	if err != nil {
		return nil, err
	}

	tags, err := entities.NewTags(cfg)
	if err != nil {
		return nil, err
	}

	themes, err := entities.NewThemes(cfg)
	if err != nil {
		return nil, err
	}

	user, err := entities.NewUser(cfg)
	if err != nil {
		return nil, err
	}

	return &domain{
		articles: articles,
		authors:  authors,
		tags:     tags,
		themes:   themes,
		user:     user,
	}, nil
}
