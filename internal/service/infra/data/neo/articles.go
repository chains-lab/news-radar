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

	CreateHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error
	DeleteHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error

	CreateAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error
	DeleteAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error

	CreateAuthorshipRelationship(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error
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
				MERGE (a)-[:HAS_TAG]->(t)
			)
			FOREACH (themeName IN $themes |
				MATCH (th:Theme { name: themeName })
				MERGE (a)-[:ABOUT]->(th)
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

func (a *articles) CreateHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (t:Tag { name: $tagName })
			MERGE (art)-[r:HAS_TAG]->(t)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tagName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create HAS_TAG relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) DeleteHasTagRelationship(ctx context.Context, articleID uuid.UUID, tagName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:HAS_TAG]->(t:Tag { name: $tagName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"tagName":   tagName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete HAS_TAG relationship: %w", err)
		}

		return nil, nil
	})

	return err
}

func (a *articles) CreateAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })
			MATCH (th:Theme { name: $themeName })
			MERGE (art)-[r:ABOUT]->(th)
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": themeName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create ABOUT relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) DeleteAboutRelationship(ctx context.Context, articleID uuid.UUID, themeName string) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (art:Article { id: $articleID })-[r:ABOUT]->(th:Theme { name: $themeName })
			DELETE r
		`
		params := map[string]any{
			"articleID": articleID.String(),
			"themeName": themeName,
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete ABOUT relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (a *articles) CreateAuthorshipRelationship(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
            MATCH (art:Article { id: $articleID })
            MATCH (auth:Author { id: $authorID })
            MERGE (art)-[:AUTHORED_BY]->(auth)
        `
		params := map[string]any{
			"articleID": articleID.String(),
			"authorID":  authorID.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create authorship relationship: %w", err)
		}
		return nil, nil
	})
	return err
}
