package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"author_id": validation.NewError("author_id", "invalid UUID format"),
		})...)
		return
	}

	author, err := h.app.GetAuthorByID(r.Context(), authorID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("Error getting author")
		return
	}

	httpkit.Render(w, responses.Author(author))
}
