package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/sirupsen/logrus"
)

type App interface {
	CreateArticle(ctx context.Context, request app.CreateArticleRequest) (models.Article, error)
	UpdateArticle(ctx context.Context, articleID uuid.UUID, request app.UpdateArticleRequest) (models.Article, error)
	DeleteArticle(ctx context.Context, articleID uuid.UUID) error
	GetArticleByID(ctx context.Context, userID, articleID uuid.UUID) (models.Article, bool, error)
	SetTags(ctx context.Context, articleID uuid.UUID, tags []string) error
	SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error
}

type Handler struct {
	log *logrus.Logger
	app App
}

func NewHandlers(log *logrus.Logger, app *app.App) *Handler {
	return &Handler{
		log: log,
		app: app,
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
}
