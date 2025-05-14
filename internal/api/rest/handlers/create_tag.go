package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/enums"
)

func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
	}

	req, err := requests.CreateTag(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	tagStatus, ok := enums.ParseTagStatus(req.Data.Attributes.Status)
	if !ok {
		h.log.WithError(err).Warn("Error parsing tag status")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Invalid tag status",
		})...)
		return
	}

	tagType, ok := enums.ParseTagType(req.Data.Attributes.Type)
	if !ok {
		h.log.WithError(err).Warn("Error parsing tag type")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Invalid tag type",
		})...)
		return
	}

	tag, err := h.app.CreateTag(r.Context(), app.CreateTagRequest{
		Name:   req.Data.Attributes.Name,
		Type:   tagType,
		Status: tagStatus,
		Color:  req.Data.Attributes.Color,
		Icon:   req.Data.Attributes.Icon,
	})

	if err != nil {
		switch {
		case errors.Is(err, ape.ErrorTagNameAlreadyTaken):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusConflict,
				Title:  "Tag name already taken",
				Detail: "Tag name already taken",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Failed to create tag")
		return
	}

	h.log.Infof("Created tag %s, by: %s", req.Data.Attributes.Name, user.AccountID)

	httpkit.Render(w, responses.Tag(tag))
}
