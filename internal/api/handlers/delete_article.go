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

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Error("Invalid article ID")
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	err = h.app.DeleteArticle(r.Context(), articleID)
	if err != nil {
		switch {
		//TODO If the article is associated with entities, it will not be deleted, we need cath this error in future.
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("Failed to delete article")
		return
	}

	h.log.Infof("Deleted article: %s by user: %s", articleID.String(), user.AccountID.String())

	w.WriteHeader(http.StatusNoContent)
}
