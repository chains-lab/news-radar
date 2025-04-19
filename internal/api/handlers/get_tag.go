package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/responses"
)

func (h *Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	tagName := chi.URLParam(r, "tag")
	res, err := h.app.GetTag(r.Context(), tagName)
	if err != nil {
		h.log.WithError(err).Error("Error getting tag")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.Tag(res))
}
