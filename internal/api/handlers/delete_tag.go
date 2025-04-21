package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	err := h.app.DeleteTag(r.Context(), tag)
	if err != nil {
		switch {
		//TODO If the tag is associated with entities, it will not be deleted.
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, problems.NotFound("tag not found"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error deleting all tag from article %s", tag)
		return
	}

	httpkit.Render(w, http.StatusNoContent)
}
