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

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	req, err := requests.UpdateAuthor(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	update := app.UpdateAuthorRequest{}
	if req.Data.Attributes.Name != nil {
		update.Name = req.Data.Attributes.Name
	}

	if req.Data.Attributes.Status != nil {
		status, ok := enums.ParseAuthorStatus(*req.Data.Attributes.Status)
		if !ok {
			h.log.Warn("Error parsing status")
			httpkit.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		update.Status = &status
	}

	if req.Data.Attributes.Desc != nil {
		update.Desc = req.Data.Attributes.Desc
	}

	if req.Data.Attributes.Avatar != nil {
		update.Avatar = req.Data.Attributes.Avatar
	}

	if req.Data.Attributes.Telegram != nil {
		update.Telegram = req.Data.Attributes.Telegram
	}

	if req.Data.Attributes.Twitter != nil {
		update.Twitter = req.Data.Attributes.Twitter
	}

	if req.Data.Attributes.Email != nil {
		update.Email = req.Data.Attributes.Email
	}

	err = h.app.UpdateAuthor(r.Context(), authorID, update)
	if err != nil {
		switch {
		case err == nil:
			h.log.WithError(err).Errorf("author id: %s", authorID)
			httpkit.RenderErr(w, problems.NotFound("author not found"))
			return
		default:
			h.log.WithError(err).Errorf("error updating author id: %s", authorID)
			httpkit.RenderErr(w, problems.InternalError())
			return
		}
	}

	author, err := h.app.GetAuthorByID(r.Context(), authorID)
	if err != nil {
		h.log.WithError(err).Error("Error getting author")
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.Author(author))
}
