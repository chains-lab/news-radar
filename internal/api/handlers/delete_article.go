package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Error("Invalid article ID")
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	err = h.app.DeleteArticle(r.Context(), articleID)
	if err != nil {
		h.log.WithError(err).Error("Failed to delete article")
		http.Error(w, "Failed to delete article", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
