package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/entities"
	"github.com/hs-zavet/news-radar/internal/config"
)

type Article interface {
}

type Author interface {
}

type Tag interface {
}

type User interface {
}

type Reaction interface {
}

type App struct {
	articles Article
	authors  Author
	tags     Tag
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

	user, err := entities.NewUser(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		articles: articles,
		authors:  authors,
		tags:     tags,
		user:     user,
	}, nil
}

func (a *App) Testacion() {
}

func (a *App) CreateArticle(Id uuid.UUID, Title string, Content string, Author string, CreatedAt time.Time, UpdatedAt time.Time) {

}
