package neodb

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type RepostsImpl struct {
	driver neo4j.Driver
}

func NewReposts(uri, username, password string) (*RepostsImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &RepostsImpl{
		driver: driver,
	}, nil
}

func (r *RepostsImpl) Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:UserModels { id: $userID })
			MATCH (a:ArticleModel { id: $articleID })
			MERGE (u)-[:REPOSTED]->(a)
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

	return nil
}

func (r *RepostsImpl) Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:User { id: $userID })-[l:REPOSTED]->(a:Article { id: $articleID })
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

func (r *RepostsImpl) GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:user { id: $userID })-[:REPOSTED]->(a:article)
			RETURN a.id
		`
		params := map[string]interface{}{
			"userID": userID.String(),
		}
		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var ids []uuid.UUID
		for cursor.Next() {
			record := cursor.Record()
			id, ok := record.Get("a.id")
			if !ok {
				continue
			}

			ids = append(ids, uuid.MustParse(id.(string)))
		}

		return ids, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]uuid.UUID), nil
}

func (r *RepostsImpl) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error) {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := `
			MATCH (u:users)-[:REPOSTED]->(a:ArticleModel { id: $articleID })
			RETURN u.id
		`
		params := map[string]interface{}{
			"articleID": articleID.String(),
		}
		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		var ids []uuid.UUID
		for cursor.Next() {
			record := cursor.Record()
			id, ok := record.Get("u.id")
			if !ok {
				continue
			}

			ids = append(ids, uuid.MustParse(id.(string)))
		}

		return ids, nil
	})
	if err != nil {
		return nil, err
	}

	return result.([]uuid.UUID), nil
}
