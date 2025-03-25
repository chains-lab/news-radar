package neodb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Likes interface {
	Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error

	GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error)
}

type likes struct {
	driver neo4j.Driver
}

func NewLikes(uri, username, password string) (Likes, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &likes{
		driver: driver,
	}, nil
}

func (l *likes) Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:User { id: $userID })
			MATCH (a:Article { id: $articleID })
			MERGE (u)-[:LIKED]->(a)
		`
		params := map[string]interface{}{
			"userID":    userID.String(),
			"articleID": articleID.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (l *likes) Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:User { id: $userID })-[l:LIKED]->(a:Article { id: $articleID })
			DELETE l
		`
		params := map[string]interface{}{
			"userID":    userID.String(),
			"articleID": articleID.String(),
		}
		_, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (l *likes) GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (u:User { id: $userID })-[l:LIKED]->(a:Article)
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
			// Предполагается, что articleID хранится как строка
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

func (l *likes) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error) {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
			MATCH (a:Article { id: $articleID })<-[l:LIKED]-(u:User)
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
