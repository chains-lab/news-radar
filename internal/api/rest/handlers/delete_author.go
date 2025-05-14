package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//TODO in future maybe need to dont use this handlers

func (h *Handler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})...)
		return
	}

	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Author ID must be a valid UUID.",
			Parametr: "author_id",
		})...)
		return
	}

	err = h.app.DeleteAuthor(r.Context(), authorID)
	if err != nil {
		switch {
		//TODO If the author is associated with entities, it will not be deleted, we need cath this error in future.
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Author not found",
				Detail:   "The requested author does not exist.",
				Parametr: "author_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
				Detail: "Failed to delete author",
			})...)
		}
		h.log.WithError(err).Errorf("error deleting author %s", authorID)
		return
	}

	h.log.Infof("Author %s successfully deleted by user: %s", authorID, user.AccountID)

	httpkit.Render(w, http.StatusNoContent)
}
