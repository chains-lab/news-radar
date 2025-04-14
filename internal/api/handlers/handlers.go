package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/sirupsen/logrus"
)

type App interface {
	CreateArticle(ctx context.Context, request app.CreateArticleRequest) (models.Article, error)
	UpdateArticle(ctx context.Context, articleID uuid.UUID, request app.UpdateArticleRequest) (models.Article, error)
	DeleteArticle(ctx context.Context, articleID uuid.UUID) error
	GetArticleByID(ctx context.Context, articleID uuid.UUID) (models.Article, error)
	SetTags(ctx context.Context, articleID uuid.UUID, tags []string) error
	SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error
}

type Handler struct {
	app App
	cfg config.Config
	log *logrus.Entry
}

func NewHandlers(cfg config.Config, log *logrus.Entry, app *app.App) Handler {
	return Handler{
		app: app,
		cfg: cfg,
		log: log,
	}
}
