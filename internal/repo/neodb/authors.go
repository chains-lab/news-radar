package neodb

import (
	"context"
	"fmt"

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

func (a *AuthorsImpl) Create(ctx context.Context, input AuthorCreateInput) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
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

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to create author: %w", err)
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

func (a *AuthorsImpl) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
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
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *AuthorsImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status enums.AuthorStatus) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (au:Author { id: $id })
				SET au.status = $status
				RETURN au
			`

			params := map[string]any{
				"id":     id.String(),
				"status": string(status),
			}

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update author: %w", err)
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

func (a *AuthorsImpl) GetByID(ctx context.Context, ID uuid.UUID) (AuthorModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return AuthorModel{}, err
	}

	defer session.Close()

	type resultWrapper struct {
		author AuthorModel
		err    error
	}

	resultChan := make(chan resultWrapper, 1)

	go func() {
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

				statusStr, ok := props["status"].(string)
				if !ok {
					return nil, fmt.Errorf("invalid status type")
				}

				status, ok := enums.ParseAuthorStatus(statusStr)
				if !ok {
					return nil, fmt.Errorf("invalid status value")
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
			resultChan <- resultWrapper{err: err}
			return
		}

		author, ok := result.(AuthorModel)
		if !ok {
			resultChan <- resultWrapper{err: fmt.Errorf("invalid result type")}
			return
		}

		resultChan <- resultWrapper{author: author}
	}()

	select {
	case res := <-resultChan:
		return res.author, res.err
	case <-ctx.Done():
		return AuthorModel{}, ctx.Err()
	}
}
