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
		r.Route("/ws", func(r chi.Router) {
			r.Get("/content", a.handlers.ArticleContentWS)
		})

		r.Route("/v1", func(r chi.Router) {
			r.Route("/articles", func(r chi.Router) {
				r.Post("/", a.handlers.CreateArticle)

				r.Route("/{article_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetArticle)
					r.With(admin).Delete("/", a.handlers.DeleteArticle)
					r.With(admin).Put("/", a.handlers.UpdateArticle)

					r.Route("/ws", func(r chi.Router) {
						r.Get("/content", a.handlers.ArticleContentWS)
					})

					r.Route("/tags", func(r chi.Router) {
						r.Get("/", a.handlers.GetArticleTags)
						r.With(auth).Put("/", a.handlers.SetHashTags)
						r.With(auth).Delete("/", a.handlers.CleanArticleTags)
						r.With(auth).Patch("/{tag}", a.handlers.AddTag)
						r.With(auth).Delete("/{tag}", a.handlers.DeleteTag)
					})

					r.Route("/authors", func(r chi.Router) {
						r.Get("/", a.handlers.GetArticleAuthors)
						r.With(auth).Put("/", a.handlers.SetAuthorship)
						r.With(auth).Delete("/", a.handlers.CleanArticleAuthors)
						r.With(auth).Patch("/{author_id}", a.handlers.AddAuthor)
						r.With(auth).Delete("/{author_id}", a.handlers.DeleteAuthor)
					})
				})
			})

			r.Route("/authors", func(r chi.Router) {
				r.With(admin).Post("/create", a.handlers.CreateAuthor)

				r.Route("/{author_id}", func(r chi.Router) {
					r.Get("/", a.handlers.GetAuthor)
					r.With(admin).Put("/", a.handlers.UpdateAuthor)
					r.With(admin).Delete("/", a.handlers.DeleteAuthor)
					r.With(admin).Get("/articles", a.handlers.GetAuthorArticles)
				})
			})

			r.With(auth).Route("/tags", func(r chi.Router) {
				r.Post("/create", a.handlers.CreateTag)

				r.Route("/{tag}", func(r chi.Router) {
					r.Get("/", a.handlers.GetTag)
					r.Put("/", a.handlers.UpdateTag)
					r.Delete("/", a.handlers.DeleteTag)
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
