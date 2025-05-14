package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		})...)
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
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Article not found",
				Detail:   "The requested article does not exist.",
				Parametr: "article_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Failed to delete article")
		return
	}

	h.log.Infof("Deleted article: %s by user: %s", articleID.String(), user.AccountID.String())

	w.WriteHeader(http.StatusNoContent)
}
