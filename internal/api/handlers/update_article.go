package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/enums"
)

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	req, err := requests.UpdateArticle(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	update := app.UpdateArticleRequest{}
	if req.Data.Attributes.Title != nil {
		update.Title = req.Data.Attributes.Title
	}
	if req.Data.Attributes.Desc != nil {
		update.Desc = req.Data.Attributes.Desc
	}
	if req.Data.Attributes.Icon != nil {
		update.Icon = req.Data.Attributes.Icon
	}
	if req.Data.Attributes.Status != nil {
		status, ok := enums.ParseArticleStatus(*req.Data.Attributes.Status)
		if !ok {
			h.log.Warn("Error parsing status")
			httpkit.RenderErr(w, problems.BadRequest(err)...)
			return
		}
		update.Status = &status
	}

	article, err := h.app.UpdateArticle(r.Context(), articleID, update)
	if err != nil {
		h.log.WithError(err).Warn("Error updating article")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.Article(article))
}
