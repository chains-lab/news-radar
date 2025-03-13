package neo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Article struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Tags      []string
	Themes    []string
}

type Articles interface {
	Create(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id uuid.UUID) error

	CreateAboutRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error
	DeleteAboutRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error

	CreateTopicRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error
	DeleteTopicRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error
}

type articles struct {
	driver neo4j.Driver
}

func NewArticles(uri, username, password string) (Articles, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &articles{
		driver: driver,
	}, nil
}

func (a *articles) Create(ctx context.Context, article *Article) error {
	if len(article.Tags) > 10 {
		return fmt.Errorf("article cannot have more than 10 tags")
	}
	if len(article.Themes) > 5 {
		return fmt.Errorf("article cannot have more than 5 themes")
	}

	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			CREATE (a:Article { id: $id, created_at: $created_at })
			FOREACH (tagName IN $tags |
				MATCH (t:Tag { name: tagName })
				MERGE (a)-[:ABOUT]->(t)
			)
			FOREACH (themeName IN $themes |
				MATCH (th:Theme { name: themeName })
				MERGE (a)-[:TOPIC]->(th)
			)
			RETURN a
		`
		params := map[string]any{
			"id":         article.ID.String(),
			"created_at": article.CreatedAt.UTC().Format(time.RFC3339),
			"tags":       article.Tags,
			"themes":     article.Themes,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create article with relationships: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) Delete(ctx context.Context, id uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $id })
			DETACH DELETE a
		`
		params := map[string]any{
			"id": id.String(),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete article: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) CreateAboutRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (t:Tag { name: $tagName })
			MERGE (art)-[r:ABOUT]->(t)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tagName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create ABOUT relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) DeleteAboutRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:ABOUT]->(t:Tag { name: $tagName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tagName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete ABOUT relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) CreateTopicRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (th:Theme { name: $themeName })
			MERGE (art)-[r:TOPIC]->(th)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": themeName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create TOPIC relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) DeleteTopicRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:TOPIC]->(th:Theme { name: $themeName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": themeName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete TOPIC relationship: %w", err)
		}
		return nil, nil
	})

	return err
}
