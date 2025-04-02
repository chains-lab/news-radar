package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type ArticleModel struct {
	ID     uuid.UUID
	Status string
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

type ArticleInsertInput struct {
	ID     uuid.UUID
	Status string
}

func (a *ArticlesImpl) Create(ctx context.Context, input ArticleInsertInput) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
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
				"id":     input.ID.String(),
				"status": input.Status,
			}

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to create article with relationships: %w", err)
			}
			return nil, nil
		})
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *ArticlesImpl) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
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
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *ArticlesImpl) UpdateStatus(ctx context.Context, ID uuid.UUID, status string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
                MATCH (a:Article { id: $id })
                SET a.status = $status
                RETURN a
            `

			params := map[string]interface{}{
				"id":     ID.String(),
				"status": status,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to set status for article: %w", err)
			}
			return nil, nil
		})
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *ArticlesImpl) GetByID(ctx context.Context, ID uuid.UUID) (ArticleModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return ArticleModel{}, err
	}
	defer session.Close()

	type resultWrapper struct {
		article ArticleModel
		err     error
	}

	resultChan := make(chan resultWrapper, 1)

	go func() {
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
				node, ok := records.Record().Get("a")
				if !ok {
					return nil, fmt.Errorf("article not found")
				}

				n, ok := node.(neo4j.Node)
				if !ok {
					return nil, fmt.Errorf("invalid node type")
				}

				props := n.Props()
				article := ArticleModel{
					ID: ID,
				}
				if status, ok := props["status"].(string); ok {
					article.Status = status
				}

				return article, nil
			}

			return nil, fmt.Errorf("article not found")
		})

		if err != nil {
			resultChan <- resultWrapper{err: err}
			return
		}

		article, ok := result.(ArticleModel)
		if !ok {
			resultChan <- resultWrapper{err: fmt.Errorf("unexpected result type")}
			return
		}

		resultChan <- resultWrapper{article: article, err: nil}
	}()

	select {
	case res := <-resultChan:
		return res.article, res.err
	case <-ctx.Done():
		return ArticleModel{}, ctx.Err()
	}
}
