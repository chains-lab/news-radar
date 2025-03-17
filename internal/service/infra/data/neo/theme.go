package neo

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/service/domain/models"
)

type Themes interface {
	Create(ctx context.Context, theme *models.Theme) error
	Delete(ctx context.Context, themeName string) error

	FindByName(ctx context.Context, name string) ([]*models.Theme, error)
	FindByID(ctx context.Context, id string) (*models.Theme, error)

	GetAll(ctx context.Context) ([]*models.Theme, error)
}

type themes struct {
	driver neo4j.Driver
}

func NewThemes(uri, username, password string) (Themes, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}
	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}
	return &themes{driver: driver}, nil
}

func (t *themes) Create(ctx context.Context, theme *models.Theme) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (th:Theme {
				name: $name,
				status: $status,
				category: $category,
				created_at: $created_at
			})
			RETURN th
		`
		params := map[string]any{
			"name":       theme.Name,
			"status":     string(theme.Status),
			"category":   string(theme.Type),
			"created_at": theme.CreatedAt.UTC().Format(time.RFC3339),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create theme: %w", err)
		}
		return nil, nil
	})
	return err
}

func (t *themes) Delete(ctx context.Context, themeName string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (th:Theme { name: $name })
			DETACH DELETE th
		`
		params := map[string]any{
			"name": themeName,
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete theme: %w", err)
		}
		return nil, nil
	})
	return err
}

func (t *themes) GetAll(ctx context.Context) ([]*models.Theme, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (th:Theme)
			OPTIONAL MATCH (th)<-[r:TOPIC]-(:Article)
			WITH th, count(r) as popularity
			RETURN th ORDER BY popularity DESC
		`
		records, err := tx.Run(cypher, nil)
		if err != nil {
			return nil, err
		}
		var themesList []*models.Theme
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("th")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			theme := &models.Theme{
				Name:   props["name"].(string),
				Status: models.ThemeStatus(props["status"].(string)),
				Type:   models.ThemeType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
					theme.CreatedAt = parsedTime
				}
			}
			themesList = append(themesList, theme)
		}
		return themesList, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.Theme), nil
}

func (t *themes) FindByName(ctx context.Context, name string) ([]*models.Theme, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (th:Theme)
			WHERE toLower(th.name) CONTAINS toLower($name)
			RETURN th
		`
		params := map[string]any{"name": name}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var themesList []*models.Theme
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("th")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			theme := &models.Theme{
				Name:   props["name"].(string),
				Status: models.ThemeStatus(props["status"].(string)),
				Type:   models.ThemeType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
					theme.CreatedAt = parsedTime
				}
			}
			themesList = append(themesList, theme)
		}
		return themesList, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*models.Theme), nil
}

func (t *themes) FindByID(ctx context.Context, id string) (*models.Theme, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (th:Theme { name: $id })
			RETURN th LIMIT 1
		`
		params := map[string]any{"id": id}
		record, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		if record.Next() {
			node, ok := record.Record().Get("th")
			if !ok {
				return nil, fmt.Errorf("theme not found")
			}
			props := node.(neo4j.Node).Props()
			theme := &models.Theme{
				Name:   props["name"].(string),
				Status: models.ThemeStatus(props["status"].(string)),
				Type:   models.ThemeType(props["category"].(string)),
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
				if err == nil {
					theme.CreatedAt = parsedTime
				}
			}
			return theme, nil
		}
		return nil, fmt.Errorf("theme not found")
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.Theme), nil
}
