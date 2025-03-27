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

type ArticleModel struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Status    models.ArticleStatus
}

type ArticlesImpl struct {
	driver neo4j.Driver
}

func NewArticles(cfg config.Config) (*ArticlesImpl, error) {
	driver, err := neo4j.NewDriver(cfg.Database.Neo4j.URI, neo4j.BasicAuth(cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &ArticlesImpl{
		driver: driver,
	}, nil
}

func (a *ArticlesImpl) Create(ctx context.Context, article *ArticleModel) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (a:Article { 
				id: $id, 
				created_at: $created_at,
				status: $status
			})		
			RETURN a
		`
		params := map[string]any{
			"id":         article.ID.String(),
			"created_at": article.CreatedAt.UTC().Format(time.RFC3339),
			"status":     article.Status,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create article with relationships: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *ArticlesImpl) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })
			DETACH DELETE a
		`
		params := map[string]any{
			"id": id.String(),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete article: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *ArticlesImpl) UpdateStatus(ctx context.Context, ID uuid.UUID, status models.ArticleStatus) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })
			SET a.status = $status
			RETURN a
		`
		params := map[string]any{
			"id":     ID.String(),
			"status": string(status),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to set status for article: %w", err)
		}
		return nil, nil
	})
	return err
}

func (a *ArticlesImpl) Get(ctx context.Context, ID uuid.UUID) (*ArticleModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })
			RETURN a
			LIMIT 1
		`
		params := map[string]any{
			"id": ID.String(),
		}

		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		if records.Next() {
			record := records.Record()
			node, ok := record.Get("a")
			if !ok {
				return nil, fmt.Errorf("article not found")
			}
			n := node.(neo4j.Node)
			props := n.Props()
			article := &ArticleModel{
				ID: ID,
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err != nil {
					return nil, fmt.Errorf("failed to parse created_at: %w", err)
				}
				article.CreatedAt = parsedTime
			}
			if statusStr, ok := props["status"].(string); ok {
				article.Status = models.ArticleStatus(statusStr)
			}
			return article, nil
		}
		return nil, fmt.Errorf("article not found")
	})
	if err != nil {
		return nil, err
	}
	return result.(*ArticleModel), nil
}
