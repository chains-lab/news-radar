package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/recovery-flow/news-radar/internal/api/handlers"
	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
	"github.com/sirupsen/logrus"
)

type UseCases interface {
	Testacion()
}

type Api struct {
	cfg *config.Config
	log *logrus.Logger

	router *chi.Mux

	handlers Handlers
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

	router := gin.Default()

	reNews := router.Group("/re-news/news-radar")
	{
		v1 := reNews.Group("/v1")
		{
			articles := v1.Group("/articles")
			{
				articles.POST("/", func(c *gin.Context) {
					// Здесь необходимо реализовать обработчик запроса
					// Например, отправляем статус OK
					c.JSON(http.StatusOK, gin.H{"status": "ok"})
				})
			}
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.cfg.Log().Error("Ошибка при запуске сервера", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.cfg.Log().Error("Ошибка при остановке сервера", err)
	}
}
