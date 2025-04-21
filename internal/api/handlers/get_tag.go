package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	tagName := chi.URLParam(r, "tag")
	res, err := h.app.GetTag(r.Context(), tagName)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, problems.NotFound("tag not found"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("Error getting tag")
		return
	}

	httpkit.Render(w, responses.Tag(res))
}
