package handlers

import (
	"errors"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	data, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	req, err := requests.NewArticleCreate(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := h.app.CreateArticle(r.Context(), app.CreateArticleRequest{
		Title: req.Data.Attributes.Title,
	})
	if err != nil {
		switch {
		case errors.Is(err, nil):
			h.log.WithError(err).Error("Error creating article")
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
	}

	resp := responses.Article(res)
	if err != nil {
		h.log.WithError(err).Error("Failed to create article")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	h.log.Infof("Created article: %s by user: %s", res.ID.String(), data.AccountID.String())

	httpkit.Render(w, resp)
}
