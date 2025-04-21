package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = h.app.DeleteAuthor(r.Context(), authorID)
	if err != nil {
		switch {
		//TODO If the author is associated with entities, it will not be deleted, we need cath this error in future.
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error deleting author %s", authorID)
		return
	}

	h.log.Infof("Author %s successfully deleted by user: %s", authorID, user.AccountID)

	httpkit.Render(w, http.StatusNoContent)
}
