// relationships_test.go
package neodb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func setupNeo4jForRelationships(t *testing.T) (neo4j.Driver, func()) {
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})
	if err != nil {
		t.Fatalf("failed to create driver: %v", err)
	}

	cleanup := func() {
		session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		if err != nil {
			t.Fatalf("failed to create session for cleanup: %v", err)
		}
		defer session.Close()

		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			queries := []string{
				"MATCH (u:User) DETACH DELETE u",
				"MATCH (a:Article) DETACH DELETE a",
				"MATCH (u:UserNeo) DETACH DELETE u",
				"MATCH (a:ArticleNeo) DETACH DELETE a",
			}
			for _, q := range queries {
				if _, err := tx.Run(q, nil); err != nil {
					return nil, err
				}
			}
			return nil, nil
		})
		if err != nil {
			t.Fatalf("failed to cleanup nodes: %v", err)
		}
	}
	cleanup()
	return driver, cleanup
}

func createTestUser(t *testing.T, driver neo4j.Driver, label string, id uuid.UUID) {
	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("failed to create session for user creation: %v", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := fmt.Sprintf("CREATE (u:%s {id: $id})", label)
		params := map[string]interface{}{
			"id": id.String(),
		}
		_, err := tx.Run(cypher, params)
		return nil, err
	})
	if err != nil {
		t.Fatalf("failed to create test user with label %s: %v", label, err)
	}
}

func createTestArticle(t *testing.T, driver neo4j.Driver, label string, id uuid.UUID) {
	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("failed to create session for article creation: %v", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		cypher := fmt.Sprintf("CREATE (a:%s {id: $id})", label)
		params := map[string]interface{}{
			"id": id.String(),
		}
		_, err := tx.Run(cypher, params)
		return nil, err
	})
	if err != nil {
		t.Fatalf("failed to create test article with label %s: %v", label, err)
	}
}

func TestLikesImpl_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	driver, cleanup := setupNeo4jForRelationships(t)
	defer cleanup()

	userID := uuid.New()
	articleID := uuid.New()
	createTestUser(t, driver, "User", userID)
	createTestArticle(t, driver, "Article", articleID)

	likes, err := neodb.NewLikes("neo4j://localhost:7687", "neo4j", "password")
	if err != nil {
		t.Fatalf("failed to create LikesImpl: %v", err)
	}

	if err := likes.Create(ctx, userID, articleID); err != nil {
		t.Fatalf("failed to create LIKED relationship: %v", err)
	}

	articles, err := likes.GetForUser(ctx, userID)
	if err != nil {
		t.Fatalf("failed to get likes for user: %v", err)
	}
	found := false
	for _, a := range articles {
		if a == articleID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected articleID %v in likes for user", articleID)
	}

	users, err := likes.GetForArticle(ctx, articleID)
	if err != nil {
		t.Fatalf("failed to get likes for article: %v", err)
	}
	found = false
	for _, u := range users {
		if u == userID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected userID %v in likes for article", userID)
	}

	if err := likes.Delete(ctx, userID, articleID); err != nil {
		t.Fatalf("failed to delete LIKED relationship: %v", err)
	}

	articles, err = likes.GetForUser(ctx, userID)
	if err != nil {
		t.Fatalf("failed to get likes for user after delete: %v", err)
	}
	for _, a := range articles {
		if a == articleID {
			t.Errorf("expected LIKED relationship to be deleted, but found articleID %v", articleID)
		}
	}
}
