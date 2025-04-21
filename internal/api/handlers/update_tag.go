package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/enums"
)

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	tagID := chi.URLParam(r, "tag")
	if tagID == "" {
		h.log.Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"tagID": errors.New("tagID not found"),
		})...)
		return
	}

	req, err := requests.UpdateTag(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if tagID != req.Data.Id {
		h.log.Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"Id": errors.New("id and tagID do not match"),
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
		case errors.Is(err, nil):
			http.Error(w, "tagID not found", http.StatusNotFound)
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}

		h.log.WithError(err).Error("Error updating tagID")
		return
	}

	httpkit.Render(w, responses.Tag(tag))
}
