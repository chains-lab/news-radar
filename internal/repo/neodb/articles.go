package neodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type ArticleModel struct {
	ID          uuid.UUID
	Status      enums.ArticleStatus
	PublishedAt *time.Time
}

type ArticlesImpl struct {
	driver neo4j.Driver
}

func NewArticles(uri, username, password string) (*ArticlesImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
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

type ArticleInsertInput struct {
	ID     uuid.UUID
	Status enums.ArticleStatus
}

func (a *ArticlesImpl) Create(ctx context.Context, input ArticleInsertInput) (ArticleModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to create new session neo4j: %w", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (a:Article { 
				id: $id, 
				status: $status
			})		
			RETURN a
		`

		params := map[string]any{
			"id":     input.ID.String(),
			"status": input.Status,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return ArticleModel{}, fmt.Errorf("failed to create article with relationships: %w", err)
		}
		return ArticleModel{}, nil
	})

	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to create article: %w", err)
	}

	article := ArticleModel{
		ID:     input.ID,
		Status: input.Status,
	}

	return article, nil
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

	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	return nil
}

type ArticleUpdateInput struct {
	Status      *enums.ArticleStatus `json:"status,omitempty"`
	PublishedAt *time.Time           `json:"published_at,omitempty"`
}

// Update patches one or more fields on Article and returns the new node.
func (a *ArticlesImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	input ArticleUpdateInput,
) (ArticleModel, error) {
	// always set updated_at
	setParts := []string{}
	params := map[string]any{
		"id": id.String(),
	}

	// optional status
	if input.Status != nil {
		setParts = append(setParts, "a.status = $status")
		params["status"] = string(*input.Status)
	}
	if input.PublishedAt != nil {
		setParts = append(setParts, "a.published_at = $published_at")
		params["published_at"] = input.PublishedAt
	}

	// if only updated_at, just load current
	if len(setParts) == 0 {
		return a.GetByID(ctx, id)
	}

	// build cypher
	cypher := fmt.Sprintf(`
        MATCH (a:Article { id: $id })
        SET %s
        RETURN a
    `, strings.Join(setParts, ", "))

	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to create new session: %w", err)
	}
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return ArticleModel{}, fmt.Errorf("failed to run update: %w", err)
		}

		if !cursor.Next() {
			return ArticleModel{}, fmt.Errorf("article not found")
		}

		nodeVal, ok := cursor.Record().Get("a")
		if !ok {
			return ArticleModel{}, fmt.Errorf("article not found")
		}

		node, ok := nodeVal.(neo4j.Node)
		if !ok {
			return ArticleModel{}, fmt.Errorf("invalid node type")
		}

		props := node.Props()

		statusStr, ok := props["status"].(string)
		if !ok {
			return ArticleModel{}, fmt.Errorf("invalid status type")
		}

		st, ok := enums.ParseArticleStatus(statusStr)
		if !ok {
			return ArticleModel{}, fmt.Errorf("unknown status: %q", statusStr)
		}

		publishedAt, ok := props["published_at"].(time.Time)
		if !ok {
			return ArticleModel{}, fmt.Errorf("invalid published_at type")
		}

		model := ArticleModel{
			ID:     id,
			Status: st,
		}
		if input.PublishedAt != nil {
			model.PublishedAt = &publishedAt
		}

		return model, nil
	})

	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to update article: %w", err)
	}

	article, ok := result.(ArticleModel)
	if !ok {
		return ArticleModel{}, fmt.Errorf("unexpected result type")
	}

	return article, nil
}

func (a *ArticlesImpl) GetByID(ctx context.Context, ID uuid.UUID) (ArticleModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return ArticleModel{}, err
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
			return ArticleModel{}, err
		}

		if records.Next() {
			nodeVal, ok := records.Record().Get("a")
			if !ok {
				return ArticleModel{}, fmt.Errorf("article not found")
			}

			node, ok := nodeVal.(neo4j.Node)
			if !ok {
				return ArticleModel{}, fmt.Errorf("invalid node type")
			}

			props := node.Props()
			article := ArticleModel{
				ID: ID,
			}

			statusStr, ok := props["status"].(string)
			if !ok {
				return ArticleModel{}, fmt.Errorf("invalid status type")
			}

			st, ok := enums.ParseArticleStatus(statusStr)
			if !ok {
				return ArticleModel{}, fmt.Errorf("unknown status value: %q", statusStr)
			}

			article.Status = st

			return article, nil
		}

		return nil, fmt.Errorf("article not found")
	})

	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to update article: %w", err)
	}

	article, ok := result.(ArticleModel)
	if !ok {
		return ArticleModel{}, fmt.Errorf("unexpected result type")
	}

	return article, nil
}
