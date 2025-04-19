package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/enums"
)

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	if tag == "" {
		h.log.Warn("Error parsing request")
		http.Error(w, "tag not found", http.StatusBadRequest)
		return
	}

	req, err := requests.UpdateTag(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		http.Error(w, "bad request", http.StatusBadRequest)
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
			http.Error(w, "tag status not found", http.StatusBadRequest)
			return
		}

		update.Status = &status
	}

	if req.Data.Attributes.Type != nil {
		tagType, ok := enums.ParseTagType(*req.Data.Attributes.Type)
		if !ok {
			h.log.Warn("Error parsing request")
			http.Error(w, "tag type not found", http.StatusBadRequest)
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

	err = h.app.UpdateTag(r.Context(), tag, update)
	if err != nil {
		switch {
		case err == nil:
			http.Error(w, "tag not found", http.StatusNotFound)
		case err == nil:
			http.Error(w, "tag already exists", http.StatusConflict)
		default:
			h.log.WithError(err).Error("Error updating tag")
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	res, err := h.app.GetTag(r.Context(), tag)
	if err != nil {
		h.log.WithError(err).Error("Error getting tag")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.Tag(res))
}
