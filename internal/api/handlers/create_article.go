package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, _, err := tokens.GetAccountData(r.Context())
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
		switch err {
		case nil:
			httpkit.RenderErr(w, problems.NotFound("article not found"))
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

	h.log.Infof("Created article: %s by user: %s", res.ID.String(), accountID.String())

	httpkit.Render(w, resp)
}
