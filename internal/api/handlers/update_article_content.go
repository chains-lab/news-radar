package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
)

func (h *Handler) ArticleContentWS(w http.ResponseWriter, r *http.Request) {
	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		h.log.WithError(err).Warn("invalid article ID")
		http.Error(w, "invalid article ID", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.WithError(err).Error("ws upgrade error")
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		return
	}

	conn.SetPongHandler(func(string) error {
		err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			return err
		}
		return nil
	})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.log.WithError(err).Warn("ws read error")
			break
		}

		wsMsg, err := requests.ArticleContentUpdate(msg)
		if err != nil {
			h.log.WithError(err).Warn("invalid ws message")
			err := conn.WriteJSON(responses.ArticleContentUpdate("error", "Invalid ws message", nil))
			if err != nil {
				return
			}
			continue
		}

		section, err := requests.UpdateSection(wsMsg.Data.Attributes.Content)
		if err != nil {
			h.log.WithError(err).Warn("failed to build section")
			err := conn.WriteJSON(responses.ArticleContentUpdate("error", "Invalid section payload", nil))
			if err != nil {
				return
			}
			continue
		}

		if err := h.app.UpdateArticleContent(
			r.Context(),
			articleID,
			int(wsMsg.Data.Attributes.SectionID),
			section,
		); err != nil {
			h.log.WithError(err).Warn("failed to update article content")
			err := conn.WriteJSON(responses.ArticleContentUpdate("error", "Failed to update content", nil))
			if err != nil {
				return
			}
			continue
		}

		article, err := h.app.GetArticleByID(r.Context(), articleID)
		if err != nil {
			h.log.WithError(err).Error("failed to fetch updated article")
			err := conn.WriteJSON(responses.ArticleContentUpdate("error", "Failed to load article", nil))
			if err != nil {
				return
			}
			continue
		}

		err = conn.WriteJSON(responses.ArticleContentUpdate("success", "Content updated", &article.Content[int(wsMsg.Data.Attributes.SectionID)]))
		if err != nil {
			return
		}
	}
}
