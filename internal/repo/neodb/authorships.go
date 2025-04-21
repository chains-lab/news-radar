package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Authorship struct {
	driver neo4j.Driver
}

func NewAuthorship(uri, username, password string) (*Authorship, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &Authorship{
		driver: driver,
	}, nil
}

func (a *Authorship) Create(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
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
			return nil, fmt.Errorf("failed to create Authorship relationship: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to create Authorship relationship: %w", err)
	}

	return nil
}

func (a *Authorship) Delete(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
				MATCH (art:Article { id: $articleID })-[r:AUTHORED_BY]->(auth:Author { id: $authorID })
				DELETE r
			`
		params := map[string]any{
			"articleID": articleID.String(),
			"authorID":  authorID.String(),
		}

		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete AUTHOR relationship: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete Authorship relationship: %w", err)
	}
	return nil
}

func (a *Authorship) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
				MATCH (a:Article { id: $id })-[:AUTHORED_BY]->(au:Author)
				RETURN au.id AS authorID
			`
		params := map[string]any{
			"id": articleID.String(),
		}

		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var authorIDs []uuid.UUID
		for records.Next() {
			record := records.Record()
			authorIDVal, ok := record.Get("authorID")
			if !ok {
				continue
			}
			if idStr, ok := authorIDVal.(string); ok {
				uid, err := uuid.Parse(idStr)
				if err != nil {
					continue
				}
				authorIDs = append(authorIDs, uid)
			}
		}
		return authorIDs, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get authors for article: %w", err)
	}

	authorIDs, ok := result.([]uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	return authorIDs, nil
}

func (a *Authorship) GetForAuthor(ctx context.Context, authorID uuid.UUID) ([]uuid.UUID, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
				MATCH (au:Author { id: $id })<-[:AUTHORED_BY]-(art:Article)
				RETURN art.id AS articleID
			`
		params := map[string]any{
			"id": authorID.String(),
		}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var articleIDs []uuid.UUID
		for records.Next() {
			record := records.Record()
			articleIDVal, ok := record.Get("articleID")
			if !ok {
				continue
			}
			articleIDStr, ok := articleIDVal.(string)
			if !ok {
				continue
			}
			parsedID, err := uuid.Parse(articleIDStr)
			if err != nil {
				continue
			}
			articleIDs = append(articleIDs, parsedID)
		}
		return articleIDs, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get articles for author: %w", err)
	}

	articleIDs, ok := result.([]uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	return articleIDs, nil
}

func (a *Authorship) SetForArticle(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		deleteCypher := `
				MATCH (a:Article { id: $id })-[r:AUTHORED_BY]->(:Author)
				DELETE r
			`
		params := map[string]any{
			"id": articleID.String(),
		}
		res, err := tx.Run(deleteCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing Authorship relationships: %w", err)
		}
		if _, err = res.Consume(); err != nil {
			return nil, err
		}

		authorIDs := make([]string, len(authors))
		for i, authID := range authors {
			authorIDs[i] = authID.String()
		}

		params["authors"] = authorIDs

		createCypher := `
				MATCH (a:Article { id: $id })
				UNWIND $authors AS authorId
				MATCH (au:Author { id: authorId })
				MERGE (a)-[:AUTHORED_BY]->(au)
			`

		res2, err := tx.Run(createCypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create new Authorship relationships: %w", err)
		}

		if _, err = res2.Consume(); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to set authors for article: %w", err)
	}

	return nil
}
