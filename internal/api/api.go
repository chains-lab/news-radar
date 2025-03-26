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
)

type API struct {
	authMW interface{}
	roleMW interface{}
	cfg    *config.Config
	app    app.App
	router *chi.Mux
}

func NewAPI(cfg *config.Config, app app.App) API {
	return API{
		authMW: tokens.AuthMdl(cfg.JWT.AccessToken.SecretKey),
		roleMW: tokens.IdentityMdl(cfg.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser),
		cfg:    cfg,
		router: chi.NewRouter(),
	}
}

func (a *API) Run(ctx context.Context) {
	a.router.Use(
		httpkit.CtxMiddleWare(
			handlers.CtxLog(a.cfg.Log()),
			handlers.CtxApp(a.app),
		),
	)

	_ = tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	_ = tokens.IdentityMdl(a.cfg.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	a.router.Route("/re-news/news-radar", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

		})
	})

	server := httpkit.StartServer(ctx, a.cfg.Server.Port, a.router, a.cfg.Log())

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, a.cfg.Log())
}
