package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

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

	if chi.URLParam(r, "article_id") != req.Data.Id {
		h.log.Warn("Article ID in URL and body do not match")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"data.id": validation.NewError("id", "article ID in URL and body do not match"),
		})...)
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
			httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
				"data.attributes.status": validation.NewError("status", "invalid status"),
			})...)
			return
		}
		update.Status = &status
	}

	article, err := h.app.UpdateArticle(r.Context(), articleID, update)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("Failed to delete article")
		return
	}

	tags, err := h.app.GetArticleTags(r.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound("article not found"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error getting article %s", articleID)
		return
	}

	authors, err := h.app.GetArticleAuthors(r.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound("article not found"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error getting article %s", articleID)
		return
	}

	h.log.Infof("article %s updated by user %s", article.ID.String(), user.AccountID.String())

	httpkit.Render(w, responses.Article(article, tags, authors))
}
