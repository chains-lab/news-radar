package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/content"
)

func (h *Handler) UpdateArticleContent(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})...)
		return
	}

	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Article ID must be a valid UUID.",
			Parametr: "article_id",
		})...)
		return
	}

	req, err := requests.UpdateArticleContent(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	if chi.URLParam(r, "article_id") != req.Data.Id {
		h.log.Warn("Article ID in URL and body do not match")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Article ID in query and in body mast be the same.",
			Parametr: "article_id",
			Pointer:  "data/id",
		})...)
		return
	}

	sections := make([]content.Section, len(req.Data.Attributes.Content))
	for i, section := range req.Data.Attributes.Content {
		s, err := content.ParseContentSection(section)
		if err != nil {
			switch {
			case errors.Is(err, content.ErrInvalidSectionType):
				httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
					Status:  http.StatusBadRequest,
					Detail:  "Invalid section type",
					Pointer: "data/attributes/content",
				})...)
			case errors.Is(err, content.ErrInvalidAudioSection):
				httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
					Status:  http.StatusBadRequest,
					Detail:  "Invalid audio section",
					Pointer: "data/attributes/content",
				})...)
			case errors.Is(err, content.ErrInvalidMediaSection):
				httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
					Status:  http.StatusBadRequest,
					Detail:  "Invalid media section",
					Pointer: "data/attributes/content",
				})...)
			case errors.Is(err, content.ErrInvalidTextSection):
				httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
					Status:  http.StatusBadRequest,
					Detail:  "Invalid text section",
					Pointer: "data/attributes/content",
				})...)
			default:
				httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
					Status: http.StatusInternalServerError,
				})...)
			}

			h.log.WithError(err).Error("failed to parse content section")
			return
		}

		sections[i] = s
	}

	article, err := h.app.GetArticleByID(r.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Article not found",
				Detail:   "Article dose not exist.",
				Parametr: "article_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Errorf("error getting article %s", articleID)
		return
	}

	tags, err := h.app.GetArticleTags(r.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Article not found",
				Detail:   "Article dose not exist.",
				Parametr: "article_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Errorf("error getting article %s", articleID)
		return
	}

	authors, err := h.app.GetArticleAuthors(r.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Article not found",
				Detail:   "Article dose not exist.",
				Parametr: "article_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Errorf("error getting article %s", articleID)
		return
	}

	h.log.Infof("article %s updated by user %s", article.ID.String(), user.AccountID.String())

	httpkit.Render(w, responses.Article(article, tags, authors))
}
