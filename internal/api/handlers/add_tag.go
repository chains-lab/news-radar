package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) AddTag(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"article_id": validation.NewError("article_id", "invalid article id"),
		})...)
		return
	}

	tagID := chi.URLParam(r, "tag_id")

	err = h.app.AddArticleTag(r.Context(), articleID, tagID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, problems.NotFound("article tag not found"))
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.Conflict("tag is inactive")) //422 надо вернуть и переделать!
		case errors.Is(err, ape.ErrTagInactive):
			httpkit.RenderErr(w, problems.Conflict("tag is inactive"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error adding tag to article %s", articleID)
		return
	}

	article, err := h.app.GetArticleByID(r.Context(), articleID)
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

	h.log.Infof("tag %s added for article %s by user: %s", tagID, articleID, user.AccountID)

	httpkit.Render(w, responses.Article(article, tags, authors))
}
