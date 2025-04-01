package neodb

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type TagModels struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type TagsImpl struct {
	driver neo4j.Driver
}

func NewTags(uri, username, password string) (*TagsImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &TagsImpl{
		driver: driver,
	}, nil
}

func (t *TagsImpl) Create(ctx context.Context, tag TagModels) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				CREATE (t:Tag {
					name: $name,
					status: $status,
					created_at: $created_at
				})
				RETURN t
			`

			params := map[string]any{
				"name":   tag.Name,
				"status": tag.Status,
			}

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to create tag: %w", err)
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

func (t *TagsImpl) Delete(ctx context.Context, name string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				DETACH DELETE t
			`

			params := map[string]any{
				"name": name,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to delete tag: %w", err)
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

func (t *TagsImpl) UpdateStatus(ctx context.Context, name string, status string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.status = $status
				RETURN t
			`

			params := map[string]any{
				"name":   name,
				"status": status,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag status: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) UpdateName(ctx context.Context, name string, newName string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)
	go func() {
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag { name: $name })
				SET t.name = $newName
				RETURN t
			`

			params := map[string]any{
				"name":    name,
				"newName": newName,
			}
			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update tag name: %w", err)
			}
			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *TagsImpl) Get(ctx context.Context, name string) (TagModels, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return TagModels{}, err
	}
	defer session.Close()

	type resultWrapper struct {
		tag TagModels
		err error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag)
				WHERE toLower(t.name) CONTAINS toLower($name)
				RETURN t
			`

			params := map[string]any{
				"name": name,
			}
			cursor, err := tx.Run(cypher, params)
			if err != nil {
				return nil, err
			}
			if cursor.Next() {
				node, ok := cursor.Record().Get("t")
				if !ok {
					return nil, fmt.Errorf("failed to find tag")
				}
				n := node.(neo4j.Node)
				props := n.Props()
				tag := TagModels{
					Name:   props["name"].(string),
					Status: props["status"].(string),
				}
				return tag, nil
			}
			return TagModels{}, fmt.Errorf("failed to find tag")
		})
		if err != nil {
			resultChan <- resultWrapper{TagModels{}, err}
			return
		}
		tag, ok := result.(TagModels)
		if !ok {
			resultChan <- resultWrapper{TagModels{}, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{tag, nil}
	}()

	select {
	case res := <-resultChan:
		return res.tag, res.err
	case <-ctx.Done():
		return TagModels{}, ctx.Err()
	}
}

func (t *TagsImpl) Select(ctx context.Context) ([]TagModels, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	type resultWrapper struct {
		tags []TagModels
		err  error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (t:Tag)
				OPTIONAL MATCH (t)<-[r:ABOUT]-(:Article)
				WITH t, count(r) as popularity
				RETURN t ORDER BY popularity DESC
			`

			cursor, err := tx.Run(cypher, nil)
			if err != nil {
				return nil, err
			}
			var tagsList []TagModels
			for cursor.Next() {
				record := cursor.Record()
				node, ok := record.Get("t")
				if !ok {
					continue
				}
				n, ok := node.(neo4j.Node)
				if !ok {
					continue
				}
				props := n.Props()
				tag := TagModels{
					Name:   props["name"].(string),
					Status: props["status"].(string),
				}
				tagsList = append(tagsList, tag)
			}
			return tagsList, nil
		})
		if err != nil {
			resultChan <- resultWrapper{nil, err}
			return
		}
		tags, ok := result.([]TagModels)
		if !ok {
			resultChan <- resultWrapper{nil, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{tags, nil}
	}()

	select {
	case res := <-resultChan:
		return res.tags, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
