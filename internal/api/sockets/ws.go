package sockets

import (
	"context"
	"net/http"

	"github.com/chains-lab/news-radar/internal/app"
	"github.com/chains-lab/news-radar/internal/app/models"
	"github.com/chains-lab/news-radar/internal/config"
	"github.com/chains-lab/news-radar/internal/content"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type App interface {
	UpdateContentSection(ctx context.Context, articleID uuid.UUID, section content.Section) (models.Article, error)
	DeleteContentSection(ctx context.Context, articleID uuid.UUID, sectionID int) (models.Article, error)
}

type WebSocket struct {
	app      App
	cfg      config.Config
	log      *logrus.Entry
	upgrader websocket.Upgrader
}

func NewWebSocket(cfg config.Config, log *logrus.Entry, app *app.App) WebSocket {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			//allowedOrigin := "https://example.com"
			//return r.Header.GetTag("Origin") == allowedOrigin
			return true
		},
	}

	return WebSocket{
		app:      app,
		cfg:      cfg,
		log:      log,
		upgrader: upgrader,
	}
}
