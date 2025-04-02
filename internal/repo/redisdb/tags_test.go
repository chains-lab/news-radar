package redisdb

import (
	"context"
	"testing"
	"time"
)

func setupRedis(t *testing.T) *TagsImpl {
	addr := "localhost:6379"
	password := ""
	DB := 1
	tags := NewTags(addr, password, DB)

	ctx := context.Background()
	if err := tags.Drop(ctx); err != nil {
		t.Fatalf("failed to drop keys: %v", err)
	}

	return tags
}

func TestTagsImpl_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	tags := setupRedis(t)

	input := TagCreateInput{
		Name:  "testTag",
		Color: "blue",
		Icon:  "star",
	}

	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	tag, err := tags.Get(ctx, "testTag")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if tag.Name != "testTag" || tag.Color != "blue" || tag.Icon != "star" {
		t.Errorf("unexpected tag data, got: %+v", tag)
	}
}

func TestTagsImpl_UpdateIcon(t *testing.T) {
	ctx := context.Background()
	tags := setupRedis(t)

	input := TagCreateInput{
		Name:  "iconTag",
		Color: "green",
		Icon:  "oldIcon",
	}

	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := tags.UpdateIcon(ctx, "iconTag", "newIcon"); err != nil {
		t.Fatalf("UpdateIcon failed: %v", err)
	}

	tag, err := tags.Get(ctx, "iconTag")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if tag.Icon != "newIcon" {
		t.Errorf("expected icon 'newIcon', got '%s'", tag.Icon)
	}
}

func TestTagsImpl_UpdateColor(t *testing.T) {
	ctx := context.Background()
	tags := setupRedis(t)

	input := TagCreateInput{
		Name:  "colorTag",
		Color: "yellow",
		Icon:  "circle",
	}

	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Обновляем цвет.
	if err := tags.UpdateColor(ctx, "colorTag", "red"); err != nil {
		t.Fatalf("UpdateColor failed: %v", err)
	}

	tag, err := tags.Get(ctx, "colorTag")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if tag.Color != "red" {
		t.Errorf("expected color 'red', got '%s'", tag.Color)
	}
}

func TestTagsImpl_Delete(t *testing.T) {
	ctx := context.Background()
	tags := setupRedis(t)

	input := TagCreateInput{
		Name:  "deleteTag",
		Color: "purple",
		Icon:  "square",
	}

	if err := tags.Create(ctx, input); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := tags.Delete(ctx, "deleteTag"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	tag, err := tags.Get(ctx, "deleteTag")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if tag.Name != "" {
		t.Errorf("expected tag to be deleted, got: %+v", tag)
	}
}

func TestTagsImpl_Drop(t *testing.T) {
	ctx := context.Background()
	tags := setupRedis(t)

	inputs := []TagCreateInput{
		{Name: "tag1", Color: "blue", Icon: "icon1"},
		{Name: "tag2", Color: "green", Icon: "icon2"},
	}

	for _, inp := range inputs {
		if err := tags.Create(ctx, inp); err != nil {
			t.Fatalf("Create failed for %s: %v", inp.Name, err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	tag1, _ := tags.Get(ctx, "tag1")
	tag2, _ := tags.Get(ctx, "tag2")
	if tag1.Name == "" || tag2.Name == "" {
		t.Fatalf("expected tags to exist before drop")
	}

	if err := tags.Drop(ctx); err != nil {
		t.Fatalf("Drop failed: %v", err)
	}

	tag1, _ = tags.Get(ctx, "tag1")
	tag2, _ = tags.Get(ctx, "tag2")
	if tag1.Name != "" || tag2.Name != "" {
		t.Errorf("expected tags to be dropped, got tag1: %+v, tag2: %+v", tag1, tag2)
	}
}
