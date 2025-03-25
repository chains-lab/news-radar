package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Dislikes interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
}

type dislikes struct {
	driver neo4j.Driver
}

func NewDislikes(uri, username, password string) (Dislikes, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &dislikes{
		driver: driver,
	}, nil
}

func (d *dislikes) Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := d.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (u:User { id: $userID })
			MATCH (a:Article { id: $articleID })
			MERGE (u)-[:DISLIKED]->(a)
		`
		params := map[string]interface{}{
			"userID":    userID.String(),
			"articleID": articleID.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create DISLIKED relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (d *dislikes) Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := d.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (u:User { id: $userID })-[r:DISLIKED]->(a:Article { id: $articleID })
			DELETE r
		`
		params := map[string]interface{}{
			"userID":    userID.String(),
			"articleID": articleID.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("failed to delete DISLIKED relationship: %w", err)
		}
		return nil, nil
	})

	return err
}

func (d *dislikes) GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	session, err := d.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (u:User { id: $userID })-[r:DISLIKED]->(a:Article)
			RETURN a.id AS articleID
		`
		params := map[string]interface{}{
			"userID": userID.String(),
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
		return nil, err
	}
	return result.([]uuid.UUID), nil
}

func (d *dislikes) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error) {
	session, err := d.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $articleID })<-[r:DISLIKED]-(u:User)
			RETURN u.id AS userID
		`
		params := map[string]interface{}{
			"articleID": articleID.String(),
		}
		records, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}
		var userIDs []uuid.UUID
		for records.Next() {
			record := records.Record()
			userIDVal, ok := record.Get("userID")
			if !ok {
				continue
			}
			userIDStr, ok := userIDVal.(string)
			if !ok {
				continue
			}
			parsedID, err := uuid.Parse(userIDStr)
			if err != nil {
				continue
			}
			userIDs = append(userIDs, parsedID)
		}
		return userIDs, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]uuid.UUID), nil
}
