package handlers

import (
	"net/http"

	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log *logrus.Logger
}

func NewHandlers(log *logrus.Logger, app *app.App) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
}
