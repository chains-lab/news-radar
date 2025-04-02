package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type LikesImpl struct {
	driver neo4j.Driver
}

func NewLikes(uri, username, password string) (*LikesImpl, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		return nil, err
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, err
	}

	return &LikesImpl{
		driver: driver,
	}, nil
}

func (l *LikesImpl) Create(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
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
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *LikesImpl) Delete(ctx context.Context, userID uuid.UUID, articleID uuid.UUID) error {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
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
		resultChan <- err
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *LikesImpl) GetForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

func (l *LikesImpl) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]uuid.UUID, error) {
	session, err := l.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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
