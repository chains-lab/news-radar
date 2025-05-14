package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/news-radar/internal/api/rest/responses"
	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TODO need a test
func (h *Handler) GetArticlesForAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Article ID must be a valid UUID.",
		})...)
		return
	}

	res, err := h.app.GetArticleForAuthor(r.Context(), authorID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Author not found",
				Detail:   "Author dose not exist.",
				Parametr: "author_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Warn("error getting article for author")
		return
	}

	h.log.Infof("len %v", len(res))

	httpkit.Render(w, responses.ArticleShortsCollection(res))
}
