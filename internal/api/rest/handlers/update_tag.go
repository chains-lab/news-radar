package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/chains-lab/news-radar/internal/api/rest/requests"
	"github.com/chains-lab/news-radar/internal/api/rest/responses"
	"github.com/chains-lab/news-radar/internal/app"
	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/chains-lab/news-radar/internal/enums"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Article ID must be a valid UUID.",
		})...)
		return
	}

	tagID := chi.URLParam(r, "tag")

	req, err := requests.UpdateTag(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	if tagID != req.Data.Id {
		h.log.Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Tag ID must be a valid UUID.",
			Parametr: "tag_id",
			Pointer:  "data/id",
		})...)
		return
	}

	update := app.UpdateTagRequest{}

	if req.Data.Attributes.Name != nil {
		update.Name = req.Data.Attributes.Name
	}

	if req.Data.Attributes.Status != nil {
		status, ok := enums.ParseTagStatus(*req.Data.Attributes.Status)
		if !ok {
			h.log.Warn("Error parsing request")
			http.Error(w, "tagID status not found", http.StatusBadRequest)
			return
		}

		update.Status = &status
	}

	if req.Data.Attributes.Type != nil {
		tagType, ok := enums.ParseTagType(*req.Data.Attributes.Type)
		if !ok {
			h.log.Warn("Error parsing request")
			http.Error(w, "tagID type not found", http.StatusBadRequest)
			return
		}

		update.Type = &tagType
	}

	if req.Data.Attributes.Color != nil {
		update.Color = req.Data.Attributes.Color
	}

	if req.Data.Attributes.Icon != nil {
		update.Icon = req.Data.Attributes.Icon
	}

	tag, err := h.app.UpdateTag(r.Context(), tagID, update)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Tag not found",
				Detail:   "Tag does not exist.",
				Pointer:  "data/id",
				Parametr: "tag_id",
			})...)
		case errors.Is(err, ape.ErrorTagNameAlreadyTaken):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusConflict,
				Title:    "Tag name already taken",
				Detail:   "Tag name already taken.",
				Pointer:  "data/attributes/name",
				Parametr: "tag_name",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}

		h.log.WithError(err).Error("Error updating tagID")
		return
	}

	h.log.Infof("Tag %s successfully updated, by user: %s", tag, user.AccountID)

	httpkit.Render(w, responses.Tag(tag))
}
