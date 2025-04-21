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

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	req, err := requests.CreateAuthor(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	author, err := h.app.CreateAuthor(r.Context(), app.CreateAuthorRequest{
		Name: req.Data.Attributes.Name,
	})
	if err != nil {
		switch {
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
	}

	h.log.Infof("Created author: %s by user: %s", req.Data.Attributes.Name, user.AccountID.String())
	httpkit.Render(w, responses.Author(author))
}
