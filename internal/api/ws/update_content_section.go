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
	"github.com/sirupsen/logrus"
)

func (s *WebSocket) ArticleContentWS(w http.ResponseWriter, r *http.Request) {
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
			s.log.WithError(err).Warn("ws read error")
			break
		}

		msgType, payload, err := requests.ParseArticleContentWS(msg)
		if err != nil {
			s.log.WithError(err).Warn("invalid ws message")
			err := conn.WriteJSON(responses.ArticleContentUpdate("error", 400, "Invalid ws message", nil))
			if err != nil {
				return
			}
			continue
		}

		switch msgType {
		case resources.ContentUpdateSection:
			upd := payload.(resources.UpdateContentSection)

			section, err := content.ParseContentSection(upd.Section)
			if err != nil {
				s.log.WithError(err).Warn("failed to build section")
				err := conn.WriteJSON(responses.ArticleContentUpdate("error", 400, "Invalid section payload", nil))
				if err != nil {
					return
				}
				continue
			}

			_, err = s.app.UpdateContentSection(
				r.Context(),
				articleID,
				section,
			)
			if err != nil {
				switch {
				case errors.Is(err, ape.ErrArticleNotFound):
					err = conn.WriteJSON(responses.ArticleContentUpdate("error", 404, "Article not found", &section))
					if err != nil {
						return
					}
				default:
					err = conn.WriteJSON(responses.ArticleContentUpdate("error", 500, "Failed to update content", &section))
					if err != nil {
						return
					}
					break
				}
				s.log.WithError(err).Warn("failed to update article content")
				continue
			}

			err = conn.WriteJSON(responses.ArticleContentUpdate("success", 200, "ContentSection updated", &section))
			if err != nil {
				return
			}

		case resources.ContentDeleteSection:
			logrus.Info("case delete")
			upd := payload.(resources.DeleteContentSection)
			_, err = s.app.DeleteContentSection(
				r.Context(),
				articleID,
				int(upd.SectionId),
			)
			if err != nil {
				switch {
				case errors.Is(err, ape.ErrArticleNotFound):
					err = conn.WriteJSON(responses.ArticleContentUpdate("error", 404, "Article not found", nil))
					if err != nil {
						return
					}
				default:
					err = conn.WriteJSON(responses.ArticleContentUpdate("error", 500, "Failed to update content", nil))
					break
				}
				s.log.WithError(err).Warn("failed to update article content")
				continue
			}

			err = conn.WriteJSON(responses.ArticleContentUpdate("success", 200, "ContentSection updated", nil))
			if err != nil {
				return
			}

		default:
			err = conn.WriteJSON(responses.ArticleContentUpdate("error", 400, "Invalid ws message type", nil))
			if err != nil {
				return
			}
		}
	}
}
