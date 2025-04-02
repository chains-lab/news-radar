package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type RepostsImpl struct {
	driver neo4j.Driver
}

func NewReposts(uri, username, password string) (*RepostsImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
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

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			cypher := `
				MATCH (u:UserNeo { id: $userID })
				MATCH (a:ArticleNeo { id: $articleID })
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
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *RepostsImpl) Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
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
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *RepostsImpl) GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	type resultWrapper struct {
		ids []uuid.UUID
		err error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			cypher := `
				MATCH (u:user { id: $userID })-[:REPOSTED]->(a:article)
				RETURN a.id AS articleID
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
				idVal, ok := record.Get("articleID")
				if !ok {
					continue
				}

				idStr, ok := idVal.(string)
				if !ok {
					continue
				}

				parsedID, err := uuid.Parse(idStr)
				if err != nil {
					continue
				}
				ids = append(ids, parsedID)
			}
			return ids, nil
		})
		if err != nil {
			resultChan <- resultWrapper{nil, err}
			return
		}

		ids, ok := result.([]uuid.UUID)
		if !ok {
			resultChan <- resultWrapper{nil, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{ids, nil}
	}()

	select {
	case res := <-resultChan:
		return res.ids, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *RepostsImpl) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error) {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	type resultWrapper struct {
		ids []uuid.UUID
		err error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			cypher := `
				MATCH (u:users)-[:REPOSTED]->(a:ArticleNeo { id: $articleID })
				RETURN u.id AS userID
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
				idVal, ok := record.Get("userID")
				if !ok {
					continue
				}

				idStr, ok := idVal.(string)
				if !ok {
					continue
				}

				parsedID, err := uuid.Parse(idStr)
				if err != nil {
					continue
				}
				ids = append(ids, parsedID)
			}
			return ids, nil
		})
		if err != nil {
			resultChan <- resultWrapper{nil, err}
			return
		}

		ids, ok := result.([]uuid.UUID)
		if !ok {
			resultChan <- resultWrapper{nil, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{ids, nil}
	}()

	select {
	case res := <-resultChan:
		return res.ids, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
