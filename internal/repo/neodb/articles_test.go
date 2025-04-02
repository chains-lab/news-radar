// articles_test.go
package neodb

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func setupNeo4jArticles(t *testing.T) *ArticlesImpl {
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "password"

	articles, err := NewArticles(uri, username, password)
	if err != nil {
		t.Fatalf("failed to create ArticlesImpl: %v", err)
	}

	session, err := articles.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("failed to create session for cleanup: %v", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run("MATCH (a:Article) DETACH DELETE a", nil)
		return nil, err
	})
	if err != nil {
		t.Fatalf("failed to clean up articles: %v", err)
	}

	return articles
}

func TestArticlesImpl_CreateAndGetByID(t *testing.T) {
	articles := setupNeo4jArticles(t)
	ctx := context.Background()

	articleID := uuid.New()
	input := ArticleInsertInput{
		ID:     articleID,
		Status: "draft",
	}

	if err := articles.Create(ctx, input); err != nil {
		t.Fatalf("failed to create article: %v", err)
	}

	article, err := articles.GetByID(ctx, articleID)
	if err != nil {
		t.Fatalf("failed to get article by ID: %v", err)
	}

	if article.ID != articleID {
		t.Errorf("expected article ID %v, got %v", articleID, article.ID)
	}
	if article.Status != "draft" {
		t.Errorf("expected status 'draft', got '%s'", article.Status)
	}
}

func TestArticlesImpl_UpdateStatus(t *testing.T) {
	articles := setupNeo4jArticles(t)
	ctx := context.Background()

	articleID := uuid.New()
	input := ArticleInsertInput{
		ID:     articleID,
		Status: "draft",
	}

	if err := articles.Create(ctx, input); err != nil {
		t.Fatalf("failed to create article: %v", err)
	}

	newStatus := "published"
	if err := articles.UpdateStatus(ctx, articleID, newStatus); err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	article, err := articles.GetByID(ctx, articleID)
	if err != nil {
		t.Fatalf("failed to get article by ID after update: %v", err)
	}
	if article.Status != newStatus {
		t.Errorf("expected updated status '%s', got '%s'", newStatus, article.Status)
	}
}

func TestArticlesImpl_Delete(t *testing.T) {
	articles := setupNeo4jArticles(t)
	ctx := context.Background()

	articleID := uuid.New()
	input := ArticleInsertInput{
		ID:     articleID,
		Status: "draft",
	}

	if err := articles.Create(ctx, input); err != nil {
		t.Fatalf("failed to create article: %v", err)
	}

	_, err := articles.GetByID(ctx, articleID)
	if err != nil {
		t.Fatalf("failed to get article before delete: %v", err)
	}

	if err := articles.Delete(ctx, articleID); err != nil {
		t.Fatalf("failed to delete article: %v", err)
	}

	_, err = articles.GetByID(ctx, articleID)
	if err == nil {
		t.Fatalf("expected error when getting deleted article, got nil")
	}
}
