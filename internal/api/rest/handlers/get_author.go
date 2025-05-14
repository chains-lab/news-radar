package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Article ID must be a valid UUID.",
		})...)
		return
	}

	author, err := h.app.GetAuthorByID(r.Context(), authorID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Article not found",
				Detail:   "Article dose not exist.",
				Parametr: "article_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Error getting author")
		return
	}

	httpkit.Render(w, responses.Author(author))
}
