package neodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Hashtag struct {
	driver neo4j.Driver
}

func NewHashtag(uri, username, password string) (*Hashtag, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create neo4j driver: %w", err)
	}

	if err = driver.VerifyConnectivity(); err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &Hashtag{
		driver: driver,
	}, nil
}

func (h *Hashtag) Create(ctx context.Context, articleID uuid.UUID, tag string) error {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (art:Article { id: $articleID })
				MATCH (t:Tag { name: $tagName })
				MERGE (art)-[r:HAS_TAG]->(t)
			`

			params := map[string]any{
				"articleID": articleID.String(),
				"tagName":   tag,
			}

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to create HAS_TAG relationship: %w", err)
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

func (h *Hashtag) Delete(ctx context.Context, articleID uuid.UUID, tag string) error {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}

	defer session.Close()

	resultChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (art:Article { id: $articleID })-[r:HAS_TAG]->(t:Tag { name: $tagName })
				DELETE r
			`

			params := map[string]any{
				"articleID": articleID.String(),
				"tagName":   tag,
			}

			_, err := tx.Run(cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to delete HAS_TAG relationship: %w", err)
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

func (h *Hashtag) GetForArticle(ctx context.Context, articleID uuid.UUID) ([]string, error) {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}
	defer session.Close()

	type resultWrapper struct {
		tags []string
		err  error
	}
	resultChan := make(chan resultWrapper, 1)

	go func() {
		res, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
				MATCH (a:Article { id: $id })-[:HAS_TAG]->(t:Tag)
				RETURN t
			`

			params := map[string]any{
				"id": articleID.String(),
			}

			records, err := tx.Run(cypher, params)
			if err != nil {
				return nil, err
			}

			var tagsList []string
			for records.Next() {
				record := records.Record()
				node, ok := record.Get("t")
				if !ok {
					continue
				}

				n, ok := node.(neo4j.Node)
				if !ok {
					continue
				}

				props := n.Props()

				tag, ok := props["name"].(string)
				if !ok {
					continue
				}

				tagsList = append(tagsList, tag)
			}
			return tagsList, nil
		})
		if err != nil {
			resultChan <- resultWrapper{nil, err}
			return
		}

		tagsList, ok := res.([]string)
		if !ok {
			resultChan <- resultWrapper{nil, fmt.Errorf("unexpected result type")}
			return
		}

		resultChan <- resultWrapper{tagsList, nil}
	}()

	select {
	case res := <-resultChan:
		return res.tags, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (h *Hashtag) GetArticlesForTag(ctx context.Context, tag string) ([]uuid.UUID, error) {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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
		res, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
			cypher := `
                MATCH (t:Tag { name: $name })<-[:HAS_TAG]-(a:Article)
                RETURN a.id AS articleID
            `
			params := map[string]any{
				"name": tag,
			}

			records, err := tx.Run(cypher, params)
			if err != nil {
				return nil, err
			}

			var articleIDs []uuid.UUID
			for records.Next() {
				rec := records.Record()
				raw, ok := rec.Get("articleID")
				if !ok {
					continue
				}
				strID, ok := raw.(string)
				if !ok {
					continue
				}
				parsed, err := uuid.Parse(strID)
				if err != nil {
					continue
				}
				articleIDs = append(articleIDs, parsed)
			}
			if err = records.Err(); err != nil {
				return nil, err
			}
			return articleIDs, nil
		})
		if err != nil {
			resultChan <- resultWrapper{nil, err}
			return
		}

		ids, ok := res.([]uuid.UUID)
		if !ok {
			resultChan <- resultWrapper{nil, fmt.Errorf("unexpected result type")}
			return
		}
		resultChan <- resultWrapper{ids, nil}
	}()

	select {
	case r := <-resultChan:
		return r.ids, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (h *Hashtag) SetForArticle(ctx context.Context, articleID uuid.UUID, tags []string) error {
	session, err := h.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return err
	}
	defer session.Close()

	errChan := make(chan error, 1)

	go func() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
			// 1) delete all old HAS_TAG relationships
			deleteCypher := `
            MATCH (a:Article { id: $id })-[r:HAS_TAG]->()
            DELETE r
        `
			params := map[string]any{"id": articleID.String()}
			res, err := tx.Run(deleteCypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to delete existing HAS_TAG relationships: %w", err)
			}
			// drain the result so driver can run next query
			if _, err = res.Consume(); err != nil {
				return nil, err
			}

			// 2) create new HAS_TAG via UNWIND
			params["tags"] = tags
			createCypher := `
            MATCH (a:Article { id: $id })
            UNWIND $tags AS tagName
            MATCH (t:Tag { name: tagName })
            MERGE (a)-[:HAS_TAG]->(t)
        `
			res2, err := tx.Run(createCypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to create new HAS_TAG relationships: %w", err)
			}
			if _, err = res2.Consume(); err != nil {
				return nil, err
			}

			return nil, nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
