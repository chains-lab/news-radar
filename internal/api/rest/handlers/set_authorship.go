package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
)

func (h *Handler) SetAuthorship(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to retrieve account data")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Article ID must be a valid UUID.",
		})...)
		return
	}

	req, err := requests.SetAuthorship(r)
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	articleID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		h.log.WithError(err).Warn("error parsing article ID")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:  http.StatusBadRequest,
			Detail:  "Article ID must be a valid UUID.",
			Pointer: "data/id",
		})...)
		return
	}

	authorIDs := make([]uuid.UUID, 0)
	for _, authorID := range req.Data.Attributes.Authors {
		author, err := uuid.Parse(authorID)
		if err != nil {
			h.log.WithError(err).Warn("error parsing authors")
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusBadRequest,
				Detail:  "Author ID must be a valid UUID.",
				Pointer: "data/attributes/authors",
			})...)
			return
		}
		authorIDs = append(authorIDs, author)
	}

	err = h.app.SetAuthors(r.Context(), articleID, authorIDs)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusNotFound,
				Title:   "Article not found",
				Detail:  "Article dose not exist.",
				Pointer: "data/id",
			})...)
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusNotFound,
				Title:   "Author not found",
				Detail:  "Author dose not exist.",
				Pointer: "data/attributes/authors",
			})...)
		case errors.Is(err, ape.ErrAuthorInactive):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusForbidden,
				Title:   "Author inactive",
				Detail:  "Some authors are inactive.",
				Pointer: "data/attributes/authors",
			})...)
		case errors.Is(err, ape.ErrAuthorReplication):
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusConflict,
				Title:   "Author replication",
				Detail:  "Some authors are already assigned to this article.",
				Pointer: "data/attributes/authors",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Errorf("error setting authorship")
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

	h.log.Infof("created authors: %s for article: %s, by user: %s", req.Data.Attributes.Authors, req.Data.Id, user.AccountID.String())

	httpkit.Render(w, responses.Article(article, tags, authors))
}
