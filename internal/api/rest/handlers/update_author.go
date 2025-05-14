package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/gatekit/httpkit"
	"github.com/chains-lab/gatekit/tokens"
	"github.com/chains-lab/news-radar/internal/api/rest/requests"
	"github.com/chains-lab/news-radar/internal/api/rest/responses"
	"github.com/chains-lab/news-radar/internal/app"
	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/chains-lab/news-radar/internal/enums"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Detail: "Article ID must be a valid UUID.",
		})...)
		return
	}

	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:   http.StatusBadRequest,
			Detail:   "Article ID must be a valid UUID.",
			Parametr: "author_id",
		})...)
		return
	}

	req, err := requests.UpdateAuthor(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status: http.StatusBadRequest,
			Error:  err,
		})...)
		return
	}

	authorIdReq, err := uuid.Parse(req.Data.Id)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:  http.StatusBadRequest,
			Detail:  "Author ID must be a valid UUID.",
			Pointer: "data/id",
		})...)
		return
	}

	if authorID != authorIdReq {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
			Status:  http.StatusBadRequest,
			Detail:  "Author ID must be the same in query and in body.",
			Pointer: "data/id",
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
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:  http.StatusBadRequest,
				Detail:  "Invalid author status",
				Pointer: "data/attributes/status",
			})...)
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
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status:   http.StatusNotFound,
				Title:    "Author not found",
				Detail:   "Author does not exist.",
				Pointer:  "data/id",
				Parametr: "author_id",
			})...)
		default:
			httpkit.RenderErr(w, httpkit.ResponseError(httpkit.ResponseErrorInput{
				Status: http.StatusInternalServerError,
			})...)
		}
		h.log.WithError(err).Error("Error updating author")
		return
	}

	h.log.Infof("Author %s successfully updated by user: %s", authorID, user.AccountID)

	httpkit.Render(w, responses.Author(author))
}
