package handlers

import (
	"errors"
	"net/http"

	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
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

	req, err := requests.AuthorCreate(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = h.app.CreateAuthor(r.Context(), app.CreateAuthorRequest{
		Name:     req.Data.Attributes.Name,
		Desc:     req.Data.Attributes.Desc,
		Avatar:   req.Data.Attributes.Avatar,
		Email:    req.Data.Attributes.Email,
		Telegram: req.Data.Attributes.Telegram,
		Twitter:  req.Data.Attributes.Twitter,
	})
	if err != nil {
		switch {
		case errors.Is(err, nil):
			h.log.WithError(err).Error("Error creating author")
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
	}

	h.log.Infof("Created author: %s by user: %s", req.Data.Attributes.Name, user.AccountID.String())
}
