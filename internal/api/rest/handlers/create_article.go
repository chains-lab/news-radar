package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	req, err := requests.CreateArticle(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	article, err := h.app.CreateArticle(r.Context(), app.CreateArticleRequest{
		Title: req.Data.Attributes.Title,
	})
	if err != nil {
		switch {
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("Failed to create article")
		return
	}

	h.log.Infof("Created article: %s by user: %s", article.ID.String(), user.AccountID.String())

	httpkit.Render(w, responses.Article(article, nil, nil))
}
