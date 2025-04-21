package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) SetAuthorship(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve account data")
		httpkit.RenderErr(w, problems.Unauthorized(err.Error()))
		return
	}

	req, err := requests.SetAuthorship(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	articleID, err := uuid.Parse(req.Data.Attributes.ArticleID)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing article ID")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	authors := make([]uuid.UUID, 0)
	for _, authorID := range req.Data.Attributes.AuthorID {
		author, err := uuid.Parse(authorID)
		if err != nil {
			h.log.WithError(err).Warn("Error parsing authors")
			httpkit.RenderErr(w, problems.BadRequest(err)...)
			return
		}
		authors = append(authors, author)
	}

	err = h.app.SetAuthors(r.Context(), articleID, authors)
	if err != nil {
		switch {
		case errors.Is(err, nil):
			h.log.WithError(err).Error("Error creating authorship")
			httpkit.RenderErr(w, problems.InternalError())
			return
		default:
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	article, err := h.app.GetArticleByID(r.Context(), articleID)
	if err != nil {
		switch {
		case errors.Is(err, nil):
			h.log.WithError(err).Error("Error retrieving article")
			httpkit.RenderErr(w, problems.InternalError())
			return
		default:
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	h.log.Infof("Created authorship: %s for article: %s, by user: %s", req.Data.Attributes.AuthorID, req.Data.Attributes.ArticleID, user.AccountID.String())

	httpkit.Render(w, responses.Article(article, nil, nil))
}
