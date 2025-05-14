package handlers

import (
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app"
)

func (h *Handler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})...)
		return
	}

	req, err := requests.CreateAuthor(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	author, err := h.app.CreateAuthor(r.Context(), app.CreateAuthorRequest{
		Name: req.Data.Attributes.Name,
	})
	if err != nil {
		switch {
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Failed to create author")
		return
	}

	h.log.Infof("Created author: %s by user: %s", req.Data.Attributes.Name, user.AccountID.String())
	httpkit.Render(w, responses.Author(author))
}
