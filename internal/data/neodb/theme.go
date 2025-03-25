package neodb

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/service/models"
)

type Theme struct {
	Name      string             `json:"name"`
	Status    models.ThemeStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
}

type Themes interface {
	Create(ctx context.Context, theme Theme) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (*Theme, error)

	UpdateStatus(ctx context.Context, name string, status models.ThemeStatus) error
	//UpdateName(ctx context.Context, name string, newName string) error

	Select(ctx context.Context) ([]*Theme, error)
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

func (t *themes) Create(ctx context.Context, theme Theme) error {
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
				created_at: $created_at
			})
			RETURN th
		`
		params := map[string]any{
			"name":       theme.Name,
			"status":     string(theme.Status),
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

func (t *themes) UpdateName(ctx context.Context, name string, newName string) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Theme { name: $name })
			SET t.name = $newName
			RETURN t
		`
		params := map[string]any{
			"name":    name,
			"newName": newName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update theme status: %w", err)
		}
		return nil, nil
	})

	return err
}

func (t *themes) UpdateStatus(ctx context.Context, name string, status models.ThemeStatus) error {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Theme { name: $name })
			SET t.status = $status
			RETURN t
		`
		params := map[string]any{
			"name":   name,
			"status": string(status),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to update theme status: %w", err)
		}
		return nil, nil
	})

	return err
}

func (t *themes) Select(ctx context.Context) ([]*Theme, error) {
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
		var themesList []*Theme
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("th")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			status, err := models.ParseThemeStatus(props["status"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to parse theme status: %w", err)
			}
			theme := &Theme{
				Name:   props["name"].(string),
				Status: status,
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
	return result.([]*Theme), nil
}

func (t *themes) Get(ctx context.Context, name string) (*Theme, error) {
	session, err := t.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (t:Theme)
			WHERE toLower(t.name) CONTAINS toLower($name)
			RETURN t
		`
		params := map[string]any{"name": name}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var theme Theme
		record := records.Record()
		node, ok := record.Get("t")
		if !ok {
			return nil, fmt.Errorf("theme not found")
		}
		props := node.(neo4j.Node).Props()
		status, err := models.ParseThemeStatus(props["status"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse theme status: %w", err)
		}
		theme = Theme{
			Name:   props["name"].(string),
			Status: status,
		}
		if createdAtStr, ok := props["created_at"].(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
			if err == nil {
				theme.CreatedAt = parsedTime
			}
		}

		return theme, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*Theme), nil
}
