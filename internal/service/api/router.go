package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/news-radar/internal/service"
	"github.com/recovery-flow/news-radar/internal/service/api/handlers"
	"github.com/recovery-flow/tokens"
	"github.com/recovery-flow/tokens/identity"
)

func Run(ctx context.Context, svc *service.Service) {
	r := chi.NewRouter()

	r.Use(
		httpkit.CtxMiddleWare(
			handlers.CtxLog(svc.Log),
			handlers.CtxApp(svc.Domain),
			handlers.CtxConfig(svc.Config),
		),
	)

	_ = tokens.AuthMdl(svc.Config.JWT.AccessToken.SecretKey)
	_ = tokens.IdentityMdl(svc.Config.JWT.AccessToken.SecretKey, identity.Admin, identity.SuperUser)

}
