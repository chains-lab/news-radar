package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) UpdateArticleContent(w http.ResponseWriter, r *http.Request) {
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

	req, err := requests.UpdateArticleContent(r)
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

	sections := make([]content.Section, len(req.Data.Attributes.Content))
	for i, section := range req.Data.Attributes.Content {
		s, err := content.ParseContentSection(section)
		if err != nil {
			switch {
			case errors.Is(err, content.ErrInvalidSectionType):
				httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
					"data/attributes/content": validation.NewError("content", "invalid section type"),
				})...)
			case errors.Is(err, content.ErrInvalidAudioSection):
				httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
					"data/attributes/content": validation.NewError("content", "invalid audio section"),
				})...)
			case errors.Is(err, content.ErrInvalidMediaSection):
				httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
					"data/attributes/content": validation.NewError("content", "invalid media section"),
				})...)
			case errors.Is(err, content.ErrInvalidTextSection):
				httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
					"data/attributes/content": validation.NewError("content", "invalid text section"),
				})...)
			default:
				httpkit.RenderErr(w, problems.InternalError())
			}
			
			h.log.WithError(err).Error("failed to parse content section")
			return
		}

		sections[i] = s
	}

	article, err := h.app.UpdateArticleContent(r.Context(), articleID, sections)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("failed to update article")
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
