package handlers

import (
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})...)
	}

	req, err := requests.CreateArticle(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	article, err := h.app.CreateArticle(r.Context(), app.CreateArticleRequest{
		Title: req.Data.Attributes.Title,
	})
	if err != nil {
		switch {
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Failed to create article")
		return
	}

	h.log.Infof("Created article: %s by user: %s", article.ID.String(), user.AccountID.String())

	httpkit.Render(w, responses.Article(article, nil, nil))
}
