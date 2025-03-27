package handlers

import (
	"net/http"

	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log *logrus.Logger
	app *app.App
}

func NewHandlers(log *logrus.Logger, app *app.App) *Handler {
	return &Handler{
		log: log,
		app: app,
	}
}
func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {

}
