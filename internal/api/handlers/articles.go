package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/sirupsen/logrus"
)

type Article interface {
	CreateArticle(ctx context.Context, request app.CreateArticleRequest) (models.Article, error)
	UpdateArticle(ctx context.Context, articleID uuid.UUID, request app.UpdateArticleRequest) (models.Article, error)
}

type Handler struct {
	log *logrus.Logger
	app Article
}

func NewHandlers(log *logrus.Logger, app *app.App) *Handler {
	return &Handler{
		log: log,
		app: app,
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
}
