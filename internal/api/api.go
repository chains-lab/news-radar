package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/news-radar/internal/api/handlers"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/identity"
	"github.com/sirupsen/logrus"
)

type Api struct {
	cfg    *config.Config
	log    *logrus.Logger
	router *chi.Mux
}

func NewAPI(cfg *config.Config) Api {
	return Api{
		log:    cfg.Log(),
		cfg:    cfg,
		router: chi.NewRouter(),
	}
}

func (a *Api) Run(ctx context.Context, app *app.App) {
	_ = tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	_ = tokens.IdentityMdl(a.cfg.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	h := handlers.NewHandlers(a.cfg.Log(), app)

	a.router.Route("/hs/news-radar", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/articles", func(r chi.Router) {
				r.Get("/", h.Test)
			})
		})
	})

	server := httpkit.StartServer(ctx, a.cfg.Server.Port, a.router, a.cfg.Log())

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, a.cfg.Log())
}
