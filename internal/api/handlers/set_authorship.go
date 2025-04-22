package handlers

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) SetAuthorship(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	req, err := requests.SetAuthorship(r)
	if err != nil {
		h.log.WithError(err).Warn("error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	articleID, err := uuid.Parse(req.Data.Id)
	if err != nil {
		h.log.WithError(err).Warn("error parsing article ID")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"article_id": validation.NewError("invalid article ID", "invalid format of article ID"),
		})...)
		return
	}

	authorIDs := make([]uuid.UUID, 0)
	for _, authorID := range req.Data.Attributes.Authors {
		author, err := uuid.Parse(authorID)
		if err != nil {
			h.log.WithError(err).Warn("error parsing authors")
			httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
				"authorIDs": validation.NewError("invalid author ID", "invalid format of author ID"),
			})...)
			return
		}
		authorIDs = append(authorIDs, author)
	}

	err = h.app.SetAuthors(r.Context(), articleID, authorIDs)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrArticleNotFound):
			httpkit.RenderErr(w, problems.NotFound("article not found"))
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, problems.NotFound("author not found"))
		case errors.Is(err, ape.ErrAuthorInactive):
			httpkit.RenderErr(w, problems.Forbidden("author is inactive"))
		case errors.Is(err, ape.ErrAuthorReplication):
			httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
				"authors": validation.NewError("authors", fmt.Sprintf("author replication in request")),
			})...)
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Errorf("error setting authorship")
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

	h.log.Infof("created authors: %s for article: %s, by user: %s", req.Data.Attributes.Authors, req.Data.Id, user.AccountID.String())

	httpkit.Render(w, responses.Article(article, tags, authors))
}
