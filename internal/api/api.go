package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
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
	auth := tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	admin := tokens.IdentityMdl(a.cfg.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

	//h := handlers.NewHandlers(a.cfg.Log(), app)

	a.router.Route("/hs/news-radar", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/public", func(r chi.Router) {

				// Public endpoints group to watch articles and interact with them, and manage article for admin
				r.Route("/articles", func(r chi.Router) {
					r.Get("/", nil)
					r.With(admin).Post("/create", nil)

					r.Route("/{article_id}", func(r chi.Router) {
						r.Get("/", nil)
						r.With(admin).Put("/", nil)
						r.With(admin).Delete("/", nil)

						r.With(auth).Route("/reactions", func(r chi.Router) {
							r.Route("/like", func(r chi.Router) {
								r.Post("/", nil)
								r.Delete("/", nil)
							})
							r.Route("/save", func(r chi.Router) {
								r.Post("/", nil)
								r.Delete("/", nil)
							})
						})

						r.Route("/tags", func(r chi.Router) {
							r.Get("/", nil)
							r.With(auth).Post("/", nil)
							r.With(auth).Delete("/", nil)
						})

						r.Route("/authors", func(r chi.Router) {
							r.Get("/", nil)
							r.With(auth).Put("/", nil)
							r.With(auth).Patch("/", nil)
							r.With(auth).Delete("/", nil)
						})

						r.Route("/rec", func(r chi.Router) {
							r.Get("/", nil)
						})
					})
				})

				// Endpoint to interact with topics
				r.Route("/topic", func(r chi.Router) {
					r.Route("/{topic_id}", func(r chi.Router) {
						r.Get("/", nil)
						r.With(admin).Put("/", nil)
						r.With(admin).Delete("/", nil)
					})
				})

				// Endpoint to interact with authors
				r.Route("/authors", func(r chi.Router) {
					r.With(admin).Post("/create", nil)

					r.Route("/{author_id}", func(r chi.Router) {
						r.Get("/", nil)
						r.With(admin).Put("/", nil)
						r.With(admin).Delete("/", nil)
					})
				})

				//Full Admin endpoints group to manage tags and topics
				r.With(auth).Route("/tags", func(r chi.Router) {
					r.Post("/create", nil)

					r.Route("/{tag_id}", func(r chi.Router) {
						r.Get("/", nil)
						r.Put("/", nil)
						r.Delete("/", nil)
					})
				})
			})
		})
	})

	server := httpkit.StartServer(ctx, a.cfg.Server.Port, a.router, a.cfg.Log())

	<-ctx.Done()
	httpkit.StopServer(context.Background(), server, a.cfg.Log())
}
