package handlers

import (
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	accountID, _, _, _, _, err := tokens.GetAccountData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	res, err := h.app.CreateArticle(r.Context(), accountID)
	if err != nil {
		switch err {
		case nil:
			httpkit.RenderErr(w, problems.NotFound("article not found"))
		//TODO: handle other errors
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
	}

	res = responses.Article(*res)
	if err != nil {
		h.log.WithError(err).Error("Failed to create article")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, res)
}
