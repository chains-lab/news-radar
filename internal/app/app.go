package app

import (
	"github.com/recovery-flow/news-radar/internal/app/entities"
	"github.com/recovery-flow/news-radar/internal/config"
)

type Article interface {
}

type Author interface {
}

type Tag interface {
}

type Theme interface {
}

type User interface {
}

type App struct {
	articles Article
	authors  Author
	tags     Tag
	themes   Theme
	user     User
}

func NewApp(cfg config.Config) (*App, error) {
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

	return &App{
		articles: articles,
		authors:  authors,
		tags:     tags,
		themes:   themes,
		user:     user,
	}, nil
}
