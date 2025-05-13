package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/httpkit"
	"github.com/hs-zavet/comtools/httpkit/problems"
	"github.com/hs-zavet/news-radar/internal/api/rest/requests"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/app"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/tokens"
)

func (h *Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	user, err := tokens.GetAccountTokenData(r.Context())
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tagID := chi.URLParam(r, "tag")

	req, err := requests.UpdateTag(r)
	if err != nil {
		h.log.WithError(err).Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if tagID != req.Data.Id {
		h.log.Warn("Error parsing request")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"Id": errors.New("id and tagID do not match"),
		})...)
		return
	}

	update := app.UpdateTagRequest{}

	if req.Data.Attributes.Name != nil {
		update.Name = req.Data.Attributes.Name
	}

	if req.Data.Attributes.Status != nil {
		status, ok := enums.ParseTagStatus(*req.Data.Attributes.Status)
		if !ok {
			h.log.Warn("Error parsing request")
			http.Error(w, "tagID status not found", http.StatusBadRequest)
			return
		}

		update.Status = &status
	}

	if req.Data.Attributes.Type != nil {
		tagType, ok := enums.ParseTagType(*req.Data.Attributes.Type)
		if !ok {
			h.log.Warn("Error parsing request")
			http.Error(w, "tagID type not found", http.StatusBadRequest)
			return
		}

		update.Type = &tagType
	}

	if req.Data.Attributes.Color != nil {
		update.Color = req.Data.Attributes.Color
	}

	if req.Data.Attributes.Icon != nil {
		update.Icon = req.Data.Attributes.Icon
	}

	tag, err := h.app.UpdateTag(r.Context(), tagID, update)
	if err != nil {
		switch {
		case errors.Is(err, ape.ErrTagNotFound):
			httpkit.RenderErr(w, problems.NotFound())
		case errors.Is(err, ape.ErrorTagNameAlreadyTaken):
			httpkit.RenderErr(w, problems.Conflict("tag name already taken"))
		default:
			httpkit.RenderErr(w, problems.InternalError())
		}

		h.log.WithError(err).Error("Error updating tagID")
		return
	}

	h.log.Infof("Tag %s successfully updated, by user: %s", tag, user.AccountID)

	httpkit.Render(w, responses.Tag(tag))
}
