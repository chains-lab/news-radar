package app

import (
	"context"
	"time"

	"github.com/google/uuid"
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

type Reaction interface {
	MakeLike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error

	MakeDislike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error

	MakeRepost(ctx context.Context, userID, articleID uuid.UUID) error
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

func (a *App) Testacion() {

}

func (a *App) CreateArticle(Id uuid.UUID, Title string, Content string, Author string, CreatedAt time.Time, UpdatedAt time.Time) {

}
