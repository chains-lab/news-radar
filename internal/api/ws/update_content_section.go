package ws

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hs-zavet/news-radar/internal/api/requests"
	"github.com/hs-zavet/news-radar/internal/api/responses"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func (s *WebSocket) ArticleContentUpdate(w http.ResponseWriter, r *http.Request) {
	articleID, err := uuid.Parse(chi.URLParam(r, "article_id"))
	if err != nil {
		s.log.WithError(err).Warn("invalid article ID")
		http.Error(w, "invalid article ID", http.StatusBadRequest)
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.WithError(err).Error("ws upgrade error")
		return
	}
	defer conn.Close()

	//initialize the timeout timer
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_ = conn.WriteMessage(websocket.PingMessage, nil)
		}
	}()
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	//handlers for handling different types of requests
	handlers := map[string]func(msg []byte) error{
		//Update section or create new section
		resources.ContentUpdateSection: func(msg []byte) error {
			req, err := requests.ParseContentSectionUpdate(msg)
			if err != nil {
				return writeErr(conn, 400, "Invalid ws message")
			}

			section, err := content.ParseContentSection(req.Section)
			if err != nil {
				return writeErr(conn, 400, "Invalid section payload")
			}

			if _, err := s.app.UpdateContentSection(r.Context(), articleID, section); err != nil {
				if errors.Is(err, ape.ErrArticleNotFound) {
					return writeErr(conn, 404, "Article not found")
				}
				return writeErr(conn, 500, "Failed to update content")
			}

			return writeOK(conn, "content section updated success", &section)
		},

		//Delete section
		resources.ContentDeleteSection: func(msg []byte) error {
			req, err := requests.ParseContentSectionDelete(msg)
			if err != nil {
				return writeErr(conn, 400, "Invalid ws message")
			}

			if _, err := s.app.DeleteContentSection(r.Context(), articleID, int(req.SectionId)); err != nil {
				if errors.Is(err, ape.ErrArticleNotFound) {
					return writeErr(conn, 404, "Article not found")
				}
				return writeErr(conn, 500, "Failed to delete content")
			}

			return writeOK(conn, "content section deleted success", nil)
		},
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.log.WithError(err).Warn("ws read error")
			break
		}

		msgType, err := requests.ParseContSectionUpdateType(msg)
		if err != nil {
			writeErr(conn, 400, "Invalid ws message")
			continue
		}

		if h, ok := handlers[msgType]; ok {
			if err := h(msg); err != nil {
				s.log.WithError(err).Warn("handler error")
			}
		} else {
			writeErr(conn, 400, "Unknown message type")
		}
	}
}

func writeErr(conn *websocket.Conn, code int, msg string) error {
	resp := responses.ArticleContentUpdate("error", code, msg, nil)
	return conn.WriteJSON(resp)
}

func writeOK(conn *websocket.Conn, msg string, data *content.Section) error {
	resp := responses.ArticleContentUpdate("success", 200, msg, data)
	return conn.WriteJSON(resp)
}
