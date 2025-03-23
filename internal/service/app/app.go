package app

import (
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/app/entities"
)

type App interface {
}

type app struct {
	articles entities.Articles
	authors  entities.Authors
	tags     entities.Tags
	themes   entities.Themes
	user     entities.User
}

func NewApp(cfg config.Config) (App, error) {
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

	return &app{
		articles: articles,
		authors:  authors,
		tags:     tags,
		themes:   themes,
		user:     user,
	}, nil
}
