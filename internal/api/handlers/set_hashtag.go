package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) SetHashTags(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	req, err := requests.SetHashtag(r)
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if req.Data.Id != chi.URLParam(r, "article_id") {
		h.log.WithError(err).Warn("Article ID mismatch")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"article_id": validation.NewError("article_id", "article ID mismatch"),
		})...)
		return
	}

	articleID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		h.log.WithError(err).Warn("error parsing article ID")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = h.app.SetArticleTags(r.Context(), articleID, req.Data.Attributes.Tags)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound("article not found"))
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, problems.NotFound("tag not found"))
		case errors.Is(err, ape.ErrTagInactive):
			httpkit.RenderErr(w, problems.Conflict(fmt.Sprintf("tag status inactive, %s", req.Data.Attributes.Tags)))
		case errors.Is(err, ape.ErrTagReplication):
			httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
				"tags": validation.NewError("tags", "tag replication in request"),
			})...)
		case errors.Is(err, ape.ErrTooManyTags):
			httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
				"tags": validation.NewError("tags", "too many tags max 10"),
			})...)
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error setting article tags")
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
		h.log.WithError(err).Errorf("error get article")
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
		h.log.WithError(err).Errorf("error getting article tags")
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
		h.log.WithError(err).Errorf("error getting article authors")
		return
	}

	h.log.Infof("Created tags: %s for article: %s, by user: %s", req.Data.Attributes.Tags, req.Data.Id, user.AccountID.String())

	httpkit.Render(w, responses.Article(article, tags, authors))
}
