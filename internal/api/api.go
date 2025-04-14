package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/news-radar/internal/api/handlers"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/tokens"
	"github.com/hs-zavet/tokens/roles"
	"github.com/sirupsen/logrus"
)

type Api struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Handler

	log *logrus.Entry
	cfg config.Config
}

func NewAPI(cfg config.Config, log *logrus.Logger, app *app.App) Api {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	hands := handlers.NewHandlers(cfg, logger, app)

	return Api{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}

func (a *Api) Run(ctx context.Context, log *logrus.Logger) {
	auth := tokens.AuthMdl(a.cfg.JWT.AccessToken.SecretKey)
	admin := tokens.AccessGrant(a.cfg.JWT.AccessToken.SecretKey, roles.Admin, roles.SuperUser)

	a.router.Route("/hs/news-radar", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/articles", func(r chi.Router) {
				r.Get("/", nil)
				r.Post("/", a.handlers.CreateArticle)

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

	a.Start(ctx, log)

	<-ctx.Done()
	a.Stop(ctx, log)
}

func (a *Api) Start(ctx context.Context, log *logrus.Logger) {
	go func() {
		a.log.Infof("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *Api) Stop(ctx context.Context, log *logrus.Logger) {
	a.log.Info("Shutting down server...")
	if err := a.server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
	}
}
