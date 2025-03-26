package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/app/models"
)

type Author struct {
	ID     uuid.UUID           `json:"id"`
	Name   string              `json:"name"`
	Status models.AuthorStatus `json:"status"`
}

type Authors interface {
	Create(ctx context.Context, author *Author) error
	Delete(ctx context.Context, ID uuid.UUID) error

	UpdateName(ctx context.Context, ID uuid.UUID, name string) error
	UpdateStatus(ctx context.Context, ID uuid.UUID, status models.AuthorStatus) error

	GetByID(ctx context.Context, ID uuid.UUID) (*Author, error)
}

type authors struct {
	driver neo4j.Driver
}

func NewAuthors(uri, username, password string) (Authors, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &authors{
		driver: driver,
	}, nil
}

func (a *authors) Create(ctx context.Context, author *Author) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (au:Author { id: $id, name: $name, status: $status })
			RETURN au
		`
		params := map[string]any{
			"id":     author.ID.String(),
			"name":   author.Name,
			"status": author.Status,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create author: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *authors) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (au:Author { id: $id })
			DETACH DELETE au
		`
		params := map[string]any{
			"id": id.String(),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete author: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *authors) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (au:Author { id: $id })
			SET au.name = $name
			RETURN au
		`

		params := map[string]any{
			"id":   id.String(),
			"name": name,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update author: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *authors) UpdateStatus(ctx context.Context, id uuid.UUID, status models.AuthorStatus) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (au:Author { id: $id })
			SET au.status = $status
			RETURN au
		`

		params := map[string]any{
			"id":     id.String(),
			"status": status,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update author: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *authors) GetByID(ctx context.Context, ID uuid.UUID) (*Author, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (au:Author { id: $id })
			RETURN au
			LIMIT 1
		`
		params := map[string]any{
			"id": ID.String(),
		}
		record, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		if record.Next() {
			node, ok := record.Record().Get("au")
			if !ok {
				return nil, fmt.Errorf("author not found")
			}
			n := node.(neo4j.Node)
			props := n.Props()

			authorID, err := uuid.Parse(props["id"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse author id: %w", err)
			}
			status, err := models.ParseAuthorStatus(props["status"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse author status: %w", err)
			}

			author := Author{
				ID:     authorID,
				Name:   props["name"].(string),
				Status: status,
			}
			return author, nil
		}
		return nil, fmt.Errorf("author not found")
	})
	if err != nil {
		return nil, err
	}
	return result.(*Author), nil
}
