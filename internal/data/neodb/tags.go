package neodb

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/app/models"
)

type TagModels struct {
	Name      string           `json:"name"`
	Status    models.TagStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
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

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (t:TagModels {
				name: $name,
				status: $status,
				created_at: $created_at
			})
			RETURN t
		`
		params := map[string]any{
			"name":       tag.Name,
			"status":     string(tag.Status),
			"created_at": tag.CreatedAt.UTC().Format(time.RFC3339),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag: %w", err)
		}
		return nil, nil
	})
	return err
}

func (t *TagsImpl) Delete(ctx context.Context, name string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:TagModels { name: $name })
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
	return err
}

func (t *TagsImpl) UpdateStatus(ctx context.Context, name string, status models.TagStatus) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:TagModels { name: $name })
			SET t.status = $status
			RETURN t
		`
		params := map[string]any{
			"name":   name,
			"status": string(status),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update tag status: %w", err)
		}
		return nil, nil
	})

	return err
}

func (t *TagsImpl) UpdateName(ctx context.Context, name string, newName string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:TagModels { name: $name })
			SET t.name = $newName
			RETURN t
		`
		params := map[string]any{
			"name":    name,
			"newName": string(newName),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update tag status: %w", err)
		}
		return nil, nil
	})

	return err
}

func (t *TagsImpl) Get(ctx context.Context, name string) (*TagModels, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:TagModels)
			WHERE toLower(t.name) CONTAINS toLower($name)
			RETURN t
		`
		params := map[string]any{"name": name}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var tag models.Tag
		record := records.Record()
		node, ok := record.Get("t")
		if !ok {
			return nil, fmt.Errorf("failed to find tag")
		}
		props := node.(neo4j.Node).Props()
		status, err := models.ParseTagStatus(props["status"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse tag status: %w", err)
		}
		tag = models.Tag{
			Name:   props["name"].(string),
			Status: status,
		}
		if createdAtStr, ok := props["created_at"].(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
			if err == nil {
				tag.CreatedAt = parsedTime
			}
		}
		return tag, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*TagModels), nil
}

func (t *TagsImpl) Select(ctx context.Context) ([]TagModels, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:TagModels)
			OPTIONAL MATCH (t)<-[r:ABOUT]-(:ArticleModel)
			WITH t, count(r) as popularity
			RETURN t ORDER BY popularity DESC
		`
		records, err := tx.Run(cypher, nil)
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
			status, err := models.ParseTagStatus(props["status"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse tag status: %w", err)
			}
			tag := &models.Tag{
				Name:   props["name"].(string),
				Status: status,
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
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
	return result.([]TagModels), nil
}
