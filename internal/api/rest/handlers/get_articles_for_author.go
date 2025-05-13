package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

// TODO need a test
func (h *Handler) GetArticlesForAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"author_id": validation.NewError("author_id", "Invalid author ID"),
		})...)
		return
	}

	res, err := h.app.GetArticleForAuthor(r.Context(), authorID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound("article not found"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Warn("error getting article for author")
		return
	}

	h.log.Infof("len %v", len(res))

	httpkit.Render(w, responses.ArticleShortsCollection(res))
}
