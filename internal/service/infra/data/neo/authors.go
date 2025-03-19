package neo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Author struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Authors interface {
	Create(ctx context.Context, author *Author) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, newName string) error

	GetByID(ctx context.Context, ID uuid.UUID) (*Author, error)
	GetArticles(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
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
			author := Author{
				ID:   authorID,
				Name: props["name"].(string),
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

func (a *authors) GetArticles(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (au:Author { id: $id })<-[:AUTHORED_BY]-(art:Article)
			RETURN art.id AS articleID
		`
		params := map[string]any{
			"id": id.String(),
		}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var articleIDs []uuid.UUID
		for records.Next() {
			record := records.Record()
			articleIDVal, ok := record.Get("articleID")
			if !ok {
				continue
			}

			articleIDStr, ok := articleIDVal.(string)
			if !ok {
				continue
			}
			parsedID, err := uuid.Parse(articleIDStr)
			if err != nil {
				continue
			}
			articleIDs = append(articleIDs, parsedID)
		}
		return articleIDs, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]uuid.UUID), nil
}
