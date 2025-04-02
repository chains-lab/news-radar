// tags_test.go
package neodb

import (
	"context"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// setupNeo4jTags создаёт экземпляр TagsImpl для тестовой базы и очищает все узлы с меткой Tag.
func setupNeo4jTags(t *testing.T) *TagsImpl {
	uri := "neo4j://localhost:7687"
	username := "neo4j"
	password := "password"

	tags, err := NewTags(uri, username, password)
	if err != nil {
		t.Fatalf("failed to create TagsImpl: %v", err)
	}

	session, err := tags.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("failed to create session for cleanup: %v", err)
	}
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run("MATCH (t:Tag) DETACH DELETE t", nil)
		return nil, err
	})
	if err != nil {
		t.Fatalf("failed to clean up tags: %v", err)
	}

	return tags
}

func TestTagsImpl_CreateAndGet(t *testing.T) {
	tags := setupNeo4jTags(t)
	ctx := context.Background()

	input := TagCreateInput{
		Name:   "TestTag",
		Status: "active",
	}

	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	retrieved, err := tags.Get(ctx, "TestTag")
	if err != nil {
		t.Fatalf("failed to get tag: %v", err)
	}
	if retrieved.Name != "TestTag" || retrieved.Status != "active" {
		t.Errorf("retrieved tag does not match input, got: %+v", retrieved)
	}
}

func TestTagsImpl_UpdateStatus(t *testing.T) {
	tags := setupNeo4jTags(t)
	ctx := context.Background()

	input := TagCreateInput{
		Name:   "StatusTag",
		Status: "active",
	}
	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	newStatus := "inactive"
	if err := tags.UpdateStatus(ctx, "StatusTag", newStatus); err != nil {
		t.Fatalf("failed to update tag status: %v", err)
	}

	retrieved, err := tags.Get(ctx, "StatusTag")
	if err != nil {
		t.Fatalf("failed to get tag after status update: %v", err)
	}
	if retrieved.Status != newStatus {
		t.Errorf("expected status %s, got %s", newStatus, retrieved.Status)
	}
}

func TestTagsImpl_UpdateName(t *testing.T) {
	tags := setupNeo4jTags(t)
	ctx := context.Background()

	input := TagCreateInput{
		Name:   "OldName",
		Status: "active",
	}
	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	newName := "NewName"
	if err := tags.UpdateName(ctx, "OldName", newName); err != nil {
		t.Fatalf("failed to update tag name: %v", err)
	}

	retrieved, err := tags.Get(ctx, newName)
	if err != nil {
		t.Fatalf("failed to get tag after name update: %v", err)
	}
	if retrieved.Name != newName {
		t.Errorf("expected name %s, got %s", newName, retrieved.Name)
	}
}

func TestTagsImpl_Delete(t *testing.T) {
	tags := setupNeo4jTags(t)
	ctx := context.Background()

	input := TagCreateInput{
		Name:   "DeleteTag",
		Status: "active",
	}
	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	if _, err := tags.Get(ctx, "DeleteTag"); err != nil {
		t.Fatalf("failed to get tag before delete: %v", err)
	}

	if err := tags.Delete(ctx, "DeleteTag"); err != nil {
		t.Fatalf("failed to delete tag: %v", err)
	}

	if _, err := tags.Get(ctx, "DeleteTag"); err == nil {
		t.Fatalf("expected error when getting deleted tag, got nil")
	}
}

func TestTagsImpl_Select(t *testing.T) {
	tags := setupNeo4jTags(t)
	ctx := context.Background()

	inputs := []TagCreateInput{
		{Name: "TagA", Status: "active"},
		{Name: "TagB", Status: "inactive"},
		{Name: "TagC", Status: "active"},
	}

	for _, inp := range inputs {
		if err := tags.Create(ctx, inp); err != nil {
			t.Fatalf("failed to create tag %s: %v", inp.Name, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	results, err := tags.Select(ctx)
	if err != nil {
		t.Fatalf("failed to select tags: %v", err)
	}

	if len(results) != len(inputs) {
		t.Errorf("expected %d tags, got %d", len(inputs), len(results))
	}
}
