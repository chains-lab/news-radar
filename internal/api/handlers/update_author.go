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

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"author_id": validation.NewError("author_id", "invalid author id"),
		})...)
		return
	}

	req, err := requests.UpdateAuthor(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	authorIdReq, err := uuid.Parse(req.Data.Id)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"author_id": validation.NewError("author_id", "invalid author id"),
		})...)
		return
	}

	if authorID != authorIdReq {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"author_id": validation.NewError("author_id", "author id mismatch"),
		})...)
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

	author, err := h.app.UpdateAuthor(r.Context(), authorID, update)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrAuthorNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}
		h.log.WithError(err).Error("Error updating author")
		return
	}

	h.log.Infof("Author %s successfully updated by user: %s", authorID, user.AccountID)

	httpkit.Render(w, responses.Author(author))
}
