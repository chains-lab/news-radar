package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/responses"
)

func (h *Handler) DeleteArticleTag(w http.ResponseWriter, r *http.Request) {
	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tagID := strings.ToLower(chi.URLParam(r, "tag_id"))

	err = h.app.DeleteArticleTag(r.Context(), articleID, tagID)
	if err != nil {
		h.log.WithError(err).Errorf("error deleting tag %s from article %s", tagID, articleID)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	article, err := h.app.GetArticleByID(r.Context(), articleID)
	if err != nil {
		h.log.WithError(err).Errorf("error getting article %s", articleID)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.Article(article, nil, nil))
}
