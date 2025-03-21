package neo

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/service/models"
)

type Tags interface {
	Create(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, tagName string) error

	FindByID(ctx context.Context, id string) (*models.Tag, error)
	FindByName(ctx context.Context, name string) (*models.Tag, error)

	GetAll(ctx context.Context) ([]*models.Tag, error)
}

type tags struct {
	driver neo4j.Driver
}

func NewTags(uri, username, password string) (Tags, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &tags{
		driver: driver,
	}, nil
}

func (t *tags) Create(ctx context.Context, tag *models.Tag) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (t:Tag {
				name: $name,
				status: $status,
				category: $category,
				created_at: $created_at
			})
			RETURN t
		`
		params := map[string]any{
			"name":       tag.Name,
			"status":     string(tag.Status),
			"category":   string(tag.Type),
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

func (t *tags) Delete(ctx context.Context, tagName string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Tag { name: $name })
			DETACH DELETE t
		`
		params := map[string]any{
			"name": tagName,
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete tag: %w", err)
		}
		return nil, nil
	})
	return err
}

func (t *tags) GetAll(ctx context.Context) ([]*models.Tag, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Tag)
			OPTIONAL MATCH (t)<-[r:ABOUT]-(:Article)
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
			tag := &models.Tag{
				Name:   props["name"].(string),
				Status: models.TagStatus(props["status"].(string)),
				Type:   models.TagType(props["category"].(string)),
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
	return result.([]*models.Tag), nil
}

func (t *tags) FindByName(ctx context.Context, name string) (*models.Tag, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Tag)
			WHERE toLower(t.name) CONTAINS toLower($name)
			RETURN t
		`
		params := map[string]any{"name": name}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var tag models.Tag
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("t")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			tag = models.Tag{
				Name:   props["name"].(string),
				Status: models.TagStatus(props["status"].(string)),
				Type:   models.TagType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
					tag.CreatedAt = parsedTime
				}
			}
		}
		return tag, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.Tag), nil
}

func (t *tags) FindByID(ctx context.Context, id string) (*models.Tag, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Tag { name: $id })
			RETURN t LIMIT 1
		`
		params := map[string]any{"id": id}
		record, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		if record.Next() {
			node, ok := record.Record().Get("t")
			if !ok {
				return nil, fmt.Errorf("tag not found")
			}
			props := node.(neo4j.Node).Props()
			tag := &models.Tag{
				Name:   props["name"].(string),
				Status: models.TagStatus(props["status"].(string)),
				Type:   models.TagType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
					tag.CreatedAt = parsedTime
				}
			}
			return tag, nil
		}
		return nil, fmt.Errorf("tag not found")
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.Tag), nil
}
