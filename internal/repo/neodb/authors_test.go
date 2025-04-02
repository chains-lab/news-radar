// authors_test.go
package neodb

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func setupNeo4jAuthors(t *testing.T) *AuthorsImpl {
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password"

	authors, err := NewAuthors(uri, username, password)
	if err != nil {
		t.Fatalf("failed to create AuthorsImpl: %v", err)
	}

	session, err := authors.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("failed to create session for cleanup: %v", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run("MATCH (au:Author) DETACH DELETE au", nil)
		return nil, err
	})
	if err != nil {
		t.Fatalf("failed to clean up authors: %v", err)
	}

	return authors
}

func TestAuthorsImpl_CreateAndGetByID(t *testing.T) {
	authors := setupNeo4jAuthors(t)
	ctx := context.Background()
	authorID := uuid.New()

	input := AuthorCreateInput{
		ID:     authorID,
		Name:   "Test Author",
		Status: "active",
	}

	if err := authors.Create(ctx, input); err != nil {
		t.Fatalf("failed to create author: %v", err)
	}

	retrieved, err := authors.GetByID(ctx, authorID)
	if err != nil {
		t.Fatalf("failed to get author by ID: %v", err)
	}

	if retrieved.Name != "Test Author" || retrieved.Status != "active" {
		t.Errorf("retrieved author does not match input, got: %+v", retrieved)
	}
}

func TestAuthorsImpl_UpdateName(t *testing.T) {
	authors := setupNeo4jAuthors(t)
	ctx := context.Background()
	authorID := uuid.New()

	input := AuthorCreateInput{
		ID:     authorID,
		Name:   "Initial Name",
		Status: "active",
	}
	if err := authors.Create(ctx, input); err != nil {
		t.Fatalf("failed to create author: %v", err)
	}

	newName := "Updated Name"
	if err := authors.UpdateName(ctx, authorID, newName); err != nil {
		t.Fatalf("failed to update author name: %v", err)
	}

	retrieved, err := authors.GetByID(ctx, authorID)
	if err != nil {
		t.Fatalf("failed to get author after name update: %v", err)
	}

	if retrieved.Name != newName {
		t.Errorf("expected name %s, got %s", newName, retrieved.Name)
	}
}

func TestAuthorsImpl_UpdateStatus(t *testing.T) {
	authors := setupNeo4jAuthors(t)
	ctx := context.Background()
	authorID := uuid.New()

	input := AuthorCreateInput{
		ID:     authorID,
		Name:   "Test Author",
		Status: "active",
	}
	if err := authors.Create(ctx, input); err != nil {
		t.Fatalf("failed to create author: %v", err)
	}

	newStatus := "inactive"
	if err := authors.UpdateStatus(ctx, authorID, newStatus); err != nil {
		t.Fatalf("failed to update author status: %v", err)
	}

	retrieved, err := authors.GetByID(ctx, authorID)
	if err != nil {
		t.Fatalf("failed to get author after status update: %v", err)
	}

	if retrieved.Status != newStatus {
		t.Errorf("expected status %s, got %s", newStatus, retrieved.Status)
	}
}

func TestAuthorsImpl_Delete(t *testing.T) {
	authors := setupNeo4jAuthors(t)
	ctx := context.Background()
	authorID := uuid.New()

	input := AuthorCreateInput{
		ID:     authorID,
		Name:   "Test Author",
		Status: "active",
	}
	if err := authors.Create(ctx, input); err != nil {
		t.Fatalf("failed to create author: %v", err)
	}

	if _, err := authors.GetByID(ctx, authorID); err != nil {
		t.Fatalf("failed to get author before delete: %v", err)
	}

	if err := authors.Delete(ctx, authorID); err != nil {
		t.Fatalf("failed to delete author: %v", err)
	}

	if _, err := authors.GetByID(ctx, authorID); err == nil {
		t.Fatalf("expected error when getting deleted author, got nil")
	}
}
