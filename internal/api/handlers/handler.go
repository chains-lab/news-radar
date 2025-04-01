package handlers

import (
	"net/http"

	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/sirupsen/logrus"
)

type App interface {
	Testacion()
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
	h.app.Testacion()
}
