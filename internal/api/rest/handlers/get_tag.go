package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	tagName := chi.URLParam(r, "tag")

	res, err := h.app.GetTag(r.Context(), tagName)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Tag not found",
				Detail:   "Tag does not exist.",
				Parametr: "tag",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Error getting tag")
		return
	}

	httpkit.Render(w, responses.Tag(res))
}
