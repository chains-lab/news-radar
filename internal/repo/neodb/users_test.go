// users_test.go
package neodb

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func setupNeo4jUsers(t *testing.T) *UsersImpl {
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password"

	users, err := NewUsers(uri, username, password)
	if err != nil {
		t.Fatalf("failed to create UsersImpl: %v", err)
	}

	session, err := users.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("failed to create session for cleanup: %v", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run("MATCH (u:User) DETACH DELETE u", nil)
		return nil, err
	})
	if err != nil {
		t.Fatalf("failed to clean up users: %v", err)
	}

	return users
}

func TestUsersImpl_CreateAndGet(t *testing.T) {
	users := setupNeo4jUsers(t)
	ctx := context.Background()
	userID := uuid.New()

	input := UserCreateInput{
		ID: userID,
	}

	if err := users.Create(ctx, input); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	retrieved, err := users.Get(ctx, userID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if retrieved.ID != userID {
		t.Errorf("expected user ID %v, got %v", userID, retrieved.ID)
	}
}

func TestUsersImpl_Delete(t *testing.T) {
	users := setupNeo4jUsers(t)
	ctx := context.Background()
	userID := uuid.New()

	input := UserCreateInput{
		ID: userID,
	}

	if err := users.Create(ctx, input); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if _, err := users.Get(ctx, userID); err != nil {
		t.Fatalf("failed to get user before delete: %v", err)
	}

	if err := users.Delete(ctx, userID); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	if _, err := users.Get(ctx, userID); err == nil {
		t.Fatalf("expected error when getting deleted user, got nil")
	}
}
