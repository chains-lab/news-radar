package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tag := chi.URLParam(r, "tag")

	err = h.app.DeleteTag(r.Context(), tag)
	if err != nil {
		switch {
		//TODO If the tag is associated with entities, it will not be deleted, we need cath this error in future.
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, problems.NotFound("tag not found"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error deleting all tag from article %s", tag)
		return
	}

	h.log.Infof("Tag %s successfully deleted by user: %s", tag, user.AccountID)

	httpkit.Render(w, http.StatusNoContent)
}
