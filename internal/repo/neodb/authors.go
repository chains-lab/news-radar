package neodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type AuthorModel struct {
	ID     uuid.UUID          `json:"id"`
	Status enums.AuthorStatus `json:"status"`
}

type AuthorsImpl struct {
	driver neo4j.Driver
}

func NewAuthors(uri, username, password string) (*AuthorsImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &AuthorsImpl{
		driver: driver,
	}, nil
}

type AuthorCreateInput struct {
	ID     uuid.UUID          `json:"id"`
	Status enums.AuthorStatus `json:"status"`
}

func (a *AuthorsImpl) Create(ctx context.Context, input AuthorCreateInput) (AuthorModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return AuthorModel{}, err
	}
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
				CREATE (au:Author { 
					id: $id,
					status: $status 
				})
				RETURN au
			`

		params := map[string]any{
			"id":     input.ID.String(),
			"status": string(input.Status),
		}

		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return AuthorModel{}, fmt.Errorf("failed to run update: %w", err)
		}

		if !cursor.Next() {
			return AuthorModel{}, fmt.Errorf("author not found")
		}

		nodeVal, ok := cursor.Record().Get("au")
		if !ok {
			return AuthorModel{}, fmt.Errorf("author not found")
		}

		node, ok := nodeVal.(neo4j.Node)
		if !ok {
			return AuthorModel{}, fmt.Errorf("invalid node type")
		}

		props := node.Props()

		statusStr, ok := props["status"].(string)
		if !ok {
			return AuthorModel{}, fmt.Errorf("invalid status type")
		}

		st, ok := enums.ParseAuthorStatus(statusStr)
		if !ok {
			return AuthorModel{}, fmt.Errorf("unknown status value: %q", statusStr)
		}

		return AuthorModel{ID: input.ID, Status: st}, nil
	})

	if err != nil {
		return AuthorModel{}, err
	}

	author, ok := result.(AuthorModel)
	if !ok {
		return AuthorModel{}, fmt.Errorf("invalid result type")
	}

	return author, nil
}

func (a *AuthorsImpl) Delete(ctx context.Context, id uuid.UUID) error {
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

	if err != nil {
		return fmt.Errorf("failed to delete author: %w", err)
	}

	return nil
}

type AuthorUpdateInput struct {
	Status *enums.AuthorStatus `json:"status,omitempty"`
}

// Update applies one or more changes to an Author node and returns the updated model.
func (a *AuthorsImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	input AuthorUpdateInput,
) (AuthorModel, error) {
	setParts := []string{}
	params := map[string]any{
		"id": id.String(),
	}

	if input.Status != nil {
		if *input.Status == "" {
			return AuthorModel{}, fmt.Errorf("status cannot be empty")
		}
		setParts = append(setParts, "au.status = $status")
		params["status"] = string(*input.Status)
	}

	if len(setParts) == 0 {
		return a.GetByID(ctx, id)
	}

	cypher := fmt.Sprintf(`
        MATCH (au:Author { id: $id })
        SET  %s
        RETURN au
    `, strings.Join(setParts, ", "))

	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return AuthorModel{}, err
	}
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return AuthorModel{}, fmt.Errorf("failed to run update: %w", err)
		}

		if !cursor.Next() {
			return AuthorModel{}, fmt.Errorf("author not found")
		}

		nodeVal, ok := cursor.Record().Get("au")
		if !ok {
			return AuthorModel{}, fmt.Errorf("author not found")
		}

		node, ok := nodeVal.(neo4j.Node)
		if !ok {
			return AuthorModel{}, fmt.Errorf("invalid node type")
		}

		props := node.Props()

		statusStr, ok := props["status"].(string)
		if !ok {
			return AuthorModel{}, fmt.Errorf("invalid status type")
		}

		st, ok := enums.ParseAuthorStatus(statusStr)
		if !ok {
			return AuthorModel{}, fmt.Errorf("unknown status value: %q", statusStr)
		}

		return AuthorModel{ID: id, Status: st}, nil
	})

	if err != nil {
		return AuthorModel{}, err
	}

	author, ok := result.(AuthorModel)
	if !ok {
		return AuthorModel{}, fmt.Errorf("invalid result type")
	}

	return author, nil
}

func (a *AuthorsImpl) GetByID(ctx context.Context, ID uuid.UUID) (AuthorModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return AuthorModel{}, err
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
			return AuthorModel{}, err
		}

		if record.Next() {
			nodeVal, ok := record.Record().Get("au")
			if !ok {
				return AuthorModel{}, fmt.Errorf("author not found")
			}

			node, ok := nodeVal.(neo4j.Node)
			if !ok {
				return AuthorModel{}, fmt.Errorf("invalid node type")
			}

			props := node.Props()

			authorID, err := uuid.Parse(props["id"].(string))
			if err != nil {
				return AuthorModel{}, fmt.Errorf("failed to parse author id: %w", err)
			}

			statusStr, ok := props["status"].(string)
			if !ok {
				return AuthorModel{}, fmt.Errorf("invalid status type")
			}

			status, ok := enums.ParseAuthorStatus(statusStr)
			if !ok {
				return AuthorModel{}, fmt.Errorf("invalid status value")
			}

			author := AuthorModel{
				ID:     authorID,
				Status: status,
			}

			return author, nil
		}

		return nil, fmt.Errorf("author not found")
	})

	if err != nil {
		return AuthorModel{}, err
	}

	author, ok := result.(AuthorModel)
	if !ok {
		return AuthorModel{}, fmt.Errorf("invalid result type")
	}

	return author, nil
}
