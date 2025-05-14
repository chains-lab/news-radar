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
)

func (h *Handler) SetHashTags(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to retrieve account data")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Article ID must be a valid UUID.",
		})...)
		return
	}

	req, err := requests.SetHashtag(r)
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	if req.Data.Id != chi.URLParam(r, "article_id") {
		h.log.WithError(err).Warn("Article ID mismatch")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:  http.StatusBadRequest,
			Detail:  "Article ID must be a valid UUID.",
			Pointer: "data/id",
		})...)
		return
	}

	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Warn("error parsing article ID")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Article ID must be a valid UUID.",
			Parametr: "article_id",
			Pointer:  "data/id",
		})...)
		return
	}

	err = h.app.SetArticleTags(r.Context(), articleID, req.Data.Attributes.Tags)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Article not found",
				Detail:   "Article does not exist.",
				Parametr: "article_id",
				Pointer:  "data/id",
			})...)
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Tag not found",
				Detail:   "Tag does not exist.",
				Parametr: "tags",
				Pointer:  "data/attributes/tags",
			})...)
		case errors.Is(err, ape.ErrTagInactive):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Tag inactive",
				Detail:   "Tag is inactive.",
				Parametr: "tags",
				Pointer:  "data/attributes/tags",
			})...)
		case errors.Is(err, ape.ErrTagReplication):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusBadRequest,
				Title:    "Tag replication",
				Detail:   "Tag is already assigned to the article.",
				Parametr: "tags",
				Pointer:  "data/attributes/tags",
			})...)
		case errors.Is(err, ape.ErrTooManyTags):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusBadRequest,
				Title:   "Too many tags",
				Detail:  "You can assign up to 10 tags to an article.",
				Pointer: "data/attributes/tags",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}

		h.log.WithError(err).Errorf("error setting article tags")
		return
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

	h.log.Infof("Created tags: %s for article: %s, by user: %s", req.Data.Attributes.Tags, req.Data.Id, user.AccountID.String())

	httpkit.Render(w, responses.Article(article, tags, authors))
}
