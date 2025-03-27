package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/news-radar/internal/api/handlers"
	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
	"github.com/sirupsen/logrus"
)

type Api struct {
	cfg      *config.Config
	log      *logrus.Logger
	router   *chi.Mux
	handlers *handlers.Handler
}

func NewAPI(cfg *config.Config, app *app.App) Api {
	h := handlers.NewHandlers(cfg.Log(), app)
	return Api{
		log:      cfg.Log(),
		cfg:      cfg,
		router:   chi.NewRouter(),
		handlers: h,
	}
}

func (a *Api) Run(ctx context.Context) {
	_ = tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	_ = tokens.IdentityMdl(a.cfg.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	a.router.Route("/re-news/news-radar", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/articles", func(r chi.Router) {
				r.Get("/", a.handlers.Test)
			})
		})
	})

	server := httpkit.StartServer(ctx, a.cfg.Server.Port, a.router, a.cfg.Log())

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, a.cfg.Log())
}
