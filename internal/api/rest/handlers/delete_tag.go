package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/go-chi/chi/v5"
)

//TODO in future maybe need to dont use this handlers

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})...)
		return
	}

	tag := chi.URLParam(r, "tag")

	err = h.app.DeleteTag(r.Context(), tag)
	if err != nil {
		switch {
		//TODO If the tag is associated with entities, it will not be deleted, we need cath this error in future.
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Tag not found",
				Detail:   "The requested tag does not exist.",
				Parametr: "tag",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
				Detail: "Failed to delete tag",
			})...)
		}
		h.log.WithError(err).Errorf("error deleting all tag from article %s", tag)
		return
	}

	h.log.Infof("Tag %s successfully deleted by user: %s", tag, user.AccountID)

	httpkit.Render(w, http.StatusNoContent)
}
