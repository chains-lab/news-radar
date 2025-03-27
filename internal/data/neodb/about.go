package neodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/recovery-flow/news-radar/internal/app/models"
)

type About struct {
	driver neo4j.Driver
}

func NewAbout(uri, username, password string) (*About, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &About{
		driver: driver,
	}, nil
}

func (a *About) Create(ctx context.Context, articleID uuid.UUID, theme string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:ArticleModel { id: $articleID })
			MATCH (th:ThemeModels { name: $theme })
			MERGE (art)-[r:ABOUT]->(th)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"theme":     theme,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create ABOUT relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *About) Delete(ctx context.Context, articleID uuid.UUID, theme string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:ArticleModel { id: $articleID })-[r:ABOUT]->(th:ThemeModels { name: $themeName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": theme,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete ABOUT relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *About) SetForArticle(ctx context.Context, articleID uuid.UUID, themes []string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		deleteCypher := `
			MATCH (a:ArticleModel { id: $articleID })-[r:ABOUT]->(:ThemeModels)
			DELETE r
		`
		params := map[string]any{"articleID": articleID.String()}
		_, err := tx.Run(deleteCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing ABOUT relationships: %w", err)
		}

		createCypher := `
			MATCH (a:ArticleModel { id: $articleID })
			FOREACH (themeName IN $ThemesImpl |
				MATCH (th:ThemeModels { name: themeName })
				MERGE (a)-[:ABOUT]->(th)
			)
		`
		params["ThemesImpl"] = themes
		_, err = tx.Run(createCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create new ABOUT relationships: %w", err)
		}
		return nil, nil
	})
	return err
}

func (a *About) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]*ThemeModels, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:ArticleModel { id: $articleID })-[:ABOUT]->(th:ThemeModels)
			RETURN th
		`
		params := map[string]any{
			"articleID": articleID.String(),
		}

		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var themesList []*ThemeModels
		for records.Next() {
			record := records.Record()
			node, ok := record.Get("th")
			if !ok {
				continue
			}
			props := node.(neo4j.Node).Props()
			status, err := models.ParseThemeStatus(props["status"].(string))
			if err != nil {
				return nil, err
			}
			theme := &ThemeModels{
				Name:   props["name"].(string),
				Status: status,
			}
			if createdAtStr, ok := props["created_at"].(string); ok {
				if parsedTime, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
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
	return result.([]*ThemeModels), nil
}
