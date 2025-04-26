package neodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func (a *ArticlesImpl) RecommendByTopic(
	ctx context.Context,
	articleID uuid.UUID,
	limit int,
) ([]ArticleModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
            MATCH (a:Article { id: $id })-[:HAS_TAG]->(t:Tag)
            WHERE t.status = 'active' AND t.type = 'topic'
            WITH a, collect(t) AS topics
            MATCH (b:Article)-[:HAS_TAG]->(t2:Tag)
            WHERE b.status = 'published'
              AND t2 IN topics
              AND b.id <> $id
            WITH b, COUNT(DISTINCT t2) AS sharedTags
            ORDER BY sharedTags DESC, b.published_at DESC
            RETURN b
            LIMIT $limit
        `
		params := map[string]any{
			"id":    articleID.String(),
			"limit": limit,
		}

		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("recommend query failed: %w", err)
		}

		var out []ArticleModel
		for cursor.Next() {
			nodeVal, _ := cursor.Record().Get("b")
			node := nodeVal.(neo4j.Node)
			props := node.Props()

			idStr := props["id"].(string)
			publishedAt := props["published_at"].(time.Time)
			statusStr := props["status"].(string)

			st, ok := enums.ParseArticleStatus(statusStr)
			if !ok {
				return nil, fmt.Errorf("unknown status %q", statusStr)
			}

			uid, err := uuid.Parse(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid uuid %q: %w", idStr, err)
			}

			out = append(out, ArticleModel{
				ID:          uid,
				Status:      st,
				PublishedAt: &publishedAt,
			})
		}
		return out, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]ArticleModel), nil
}

func (a *ArticlesImpl) TopicSearch(
	ctx context.Context,
	tag string,
	start, limit int,
) ([]ArticleModel, error) {
	session, err := a.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, fmt.Errorf("failed to open session: %w", err)
	}
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		cypher := `
            MATCH (t:Tag { id: $tag })<-[:HAS_TAG]-(a:Article)
            WHERE a.status = 'published'
            RETURN a
            ORDER BY a.published_at DESC
            SKIP $start
            LIMIT $limit
        `
		params := map[string]any{
			"tag":   tag,
			"start": start,
			"limit": limit,
		}

		cursor, err := tx.Run(cypher, params)
		if err != nil {
			return nil, fmt.Errorf("topic search query failed: %w", err)
		}

		var out []ArticleModel
		for cursor.Next() {
			nodeVal, _ := cursor.Record().Get("a")
			node := nodeVal.(neo4j.Node)
			props := node.Props()

			idStr := props["id"].(string)
			publishedAt := props["published_at"].(time.Time)
			statusStr := props["status"].(string)

			st, ok := enums.ParseArticleStatus(statusStr)
			if !ok {
				return nil, fmt.Errorf("unknown status %q", statusStr)
			}

			uid, err := uuid.Parse(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid uuid %q: %w", idStr, err)
			}

			out = append(out, ArticleModel{
				ID:          uid,
				Status:      st,
				PublishedAt: &publishedAt,
			})
		}
		return out, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]ArticleModel), nil
}
