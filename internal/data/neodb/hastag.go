package neodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
)

type Hashtag struct {
	driver neo4j.Driver
}

func NewHashtag(cfg config.Config) (*Hashtag, error) {
	driver, err := neo4j.NewDriver(cfg.Database.Neo4j.URI, neo4j.BasicAuth(cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &Hashtag{
		driver: driver,
	}, nil
}

func (h *Hashtag) Create(ctx context.Context, articleID uuid.UUID, tag string) error {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (t:Tag { name: $tagName })
			MERGE (art)-[r:HAS_TAG]->(t)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tag,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create HAS_TAG relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (h *Hashtag) Delete(ctx context.Context, articleID uuid.UUID, tag string) error {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:HAS_TAG]->(t:Tag { name: $tagName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tag,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete HAS_TAG relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (h *Hashtag) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]*models.Tag, error) {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (h:Article { id: $id })-[:HAS_TAG]->(t:Tag)
			RETURN t
		`
		params := map[string]any{
			"id": articleID.String(),
		}

		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var tagsList []*models.Tag
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("t")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			tag := &models.Tag{
				Name:   props["name"].(string),
				Status: models.TagStatus(props["status"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				if parsedTime, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
					tag.CreatedAt = parsedTime
				}
			}
			tagsList = append(tagsList, tag)
		}
		return tagsList, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.Tag), nil
}

func (h *Hashtag) SetForArticle(ctx context.Context, articleID uuid.UUID, tags []string) error {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		deleteCypher := `
			MATCH (h:Article { id: $id })-[r:HAS_TAG]->(:Tag)
			DELETE r
		`
		params := map[string]any{"id": articleID.String()}
		_, err := tx.Run(deleteCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing HAS_TAG relationships: %w", err)
		}

		createCypher := `
			MATCH (h:Article { id: $id })
			FOREACH (tagName IN $tags |
				MATCH (t:Tag { name: tagName })
				MERGE (h)-[:HAS_TAG]->(t)
			)
		`
		params["$tags"] = tags
		_, err = tx.Run(createCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create new HAS_TAG relationships: %w", err)
		}
		return nil, nil
	})
	return err
}
