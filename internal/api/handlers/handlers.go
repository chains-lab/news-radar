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

	CreateTag(ctx context.Context, request app.CreateTagRequest) error
	DeleteTag(ctx context.Context, name string) error
	UpdateTag(ctx context.Context, name string, request app.UpdateTagRequest) error
	Get(ctx context.Context, name string) (models.Tag, error)

	CreateAuthor(ctx context.Context, request app.CreateAuthorRequest) error
	UpdateAuthor(ctx context.Context, authorID uuid.UUID, request app.UpdateAuthorRequest) error
	DeleteAuthor(ctx context.Context, authorID uuid.UUID) error
	GetAuthorByID(ctx context.Context, authorID uuid.UUID) (models.Author, error)
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
