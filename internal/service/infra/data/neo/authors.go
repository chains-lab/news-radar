package neo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Author — узел автора в Neo4j.
type Author struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Authors — интерфейс для работы с узлами Author.
type Authors interface {
	// Create создаёт узел Author.
	Create(ctx context.Context, author *Author) error
	// Delete удаляет узел Author по ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// Update обновляет только имя узла Author по его ID.
	Update(ctx context.Context, id uuid.UUID, newName string) error
}

type authors struct {
	driver neo4j.Driver
}

// NewAuthors создаёт новый репозиторий для работы с авторами.
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
			CREATE (au:Author { id: $id, name: $name })
			RETURN au
		`
		params := map[string]any{
			"id":   author.ID.String(),
			"name": author.Name,
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

func (a *authors) Update(ctx context.Context, id uuid.UUID, newName string) error {
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
			"name": newName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update author: %w", err)
		}
		return nil, nil
	})

	return err
}
