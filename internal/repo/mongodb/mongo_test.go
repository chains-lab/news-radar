package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/content"
	"go.mongodb.org/mongo-driver/bson"
)

func setupArticle(t *testing.T) *ArticlesQ {
	aq, err := NewArticles("test", "mongodb://localhost:7100")
	if err != nil {
		t.Fatalf("failed to create ArticlesQ: %v", err)
	}
	if err := aq.collection.Drop(context.Background()); err != nil {
		t.Fatalf("failed to drop test collection: %v", err)
	}
	aq.filters = bson.M{}
	aq.limit = 0
	aq.skip = 0
	aq.sort = bson.D{}
	return aq
}

func TestArticles_Insert_Get_Count_Delete(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	articleID := uuid.New()
	now := time.Now().UTC()
	insertInput := ArticleInsertInput{
		ID:        articleID,
		Title:     "Test Article",
		Icon:      "test.png",
		Desc:      "Description",
		Content:   []content.Section{{Section: content.SectionTypeText, Content: map[string]any{}}},
		CreatedAt: now,
	}

	if err := aq.Insert(ctx, insertInput); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	count, err := aq.Count(ctx)
	if err != nil {
		t.Fatalf("Count failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	aq.FilterID(articleID)
	article, err := aq.Get(ctx)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if article.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got %s", article.Title)
	}

	newTitle := "Updated Title"
	newLikes := 10
	updateInput := ArticleUpdateInput{
		Title:     &newTitle,
		Likes:     &newLikes,
		UpdatedAt: time.Now().UTC(),
	}
	updatedArticle, err := aq.Update(ctx, updateInput)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updatedArticle.Title != newTitle {
		t.Errorf("Expected updated title %s, got %s", newTitle, updatedArticle.Title)
	}
	if updatedArticle.Likes != newLikes {
		t.Errorf("Expected likes %d, got %d", newLikes, updatedArticle.Likes)
	}

	if err := aq.Delete(ctx); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	count, err = aq.Count(ctx)
	if err != nil {
		t.Fatalf("Count after delete failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after delete, got %d", count)
	}
}

func TestArticles_FilterTitle_Sort_Limit_Skip(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	articles := []ArticleInsertInput{
		{
			ID:        uuid.New(),
			Title:     "Alpha Article",
			Icon:      "alpha.png",
			Desc:      "Alpha description",
			Content:   []content.Section{{Section: content.SectionTypeText, Content: map[string]any{}}},
			CreatedAt: time.Now().Add(-10 * time.Minute),
		},
		{
			ID:        uuid.New(),
			Title:     "Beta Article",
			Icon:      "beta.png",
			Desc:      "Beta description",
			Content:   []content.Section{{Section: content.SectionTypeText, Content: map[string]any{}}},
			CreatedAt: time.Now().Add(-5 * time.Minute),
		},
		{
			ID:        uuid.New(),
			Title:     "Alpha Test",
			Icon:      "alpha2.png",
			Desc:      "Alpha test description",
			Content:   []content.Section{{Section: content.SectionTypeText, Content: map[string]any{}}},
			CreatedAt: time.Now(),
		},
	}

	for _, art := range articles {
		if err := aq.Insert(ctx, art); err != nil {
			t.Fatalf("failed to insert article: %v", err)
		}
	}

	aq.filters = bson.M{}
	aq.FilterTitle("Alpha")
	results, err := aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with FilterTitle failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Expected 2 articles with 'Alpha' in title, got %d", len(results))
	}

	aq.filters = bson.M{}
	aq.Sort("created_at", true)
	results, err = aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with Sort failed: %v", err)
	}
	if len(results) < 3 {
		t.Errorf("Expected at least 3 articles, got %d", len(results))
	} else {
		if results[0].CreatedAt.After(results[1].CreatedAt) {
			t.Errorf("articles not sorted ascending by CreatedAt")
		}
	}

	aq.filters = bson.M{}
	aq.Sort("created_at", true)
	aq.Limit(1)
	aq.Skip(1)
	results, err = aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with Limit/Skip failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 article after applying limit and skip, got %d", len(results))
	}
}

func TestArticles_FilterDate(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	articleID := uuid.New()
	now := time.Now().UTC()
	insertInput := ArticleInsertInput{
		ID:        articleID,
		Title:     "Date Test Article",
		Icon:      "date.png",
		Desc:      "Date test description",
		Content:   []content.Section{{Section: content.SectionTypeText, Content: map[string]any{}}},
		CreatedAt: now,
	}
	if err := aq.Insert(ctx, insertInput); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newTime := time.Now().UTC().Add(5 * time.Minute)
	bewTitle := "Updated Date Test Article Title"

	updateInput := ArticleUpdateInput{
		Title:     &bewTitle,
		UpdatedAt: newTime,
	}
	aq.filters = bson.M{}
	aq.FilterID(articleID)
	_, err := aq.Update(ctx, updateInput)
	if err != nil {
		t.Fatalf("Update for FilterDate failed: %v", err)
	}

	aq.filters = bson.M{}
	aq.FilterDate(map[string]any{"updated_at": newTime}, true)
	results, err := aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with FilterDate failed: %v", err)
	}
	if len(results) == 0 {
		t.Errorf("Expected at least 1 article with updated_at >= %v, got 0", newTime)
	}
}

func TestArticles_MultipleUpdates(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	articleID := uuid.New()
	now := time.Now().UTC()
	insertInput := ArticleInsertInput{
		ID:        articleID,
		Title:     "Original Title",
		Icon:      "orig.png",
		Desc:      "Original description",
		Content:   []content.Section{{Section: content.SectionTypeText, Content: map[string]any{}}},
		CreatedAt: now,
	}
	if err := aq.Insert(ctx, insertInput); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	newTitle := "New Title"
	updateInput1 := ArticleUpdateInput{
		Title:     &newTitle,
		UpdatedAt: time.Now().UTC(),
	}
	aq.filters = bson.M{}
	aq.FilterID(articleID)
	updatedArticle, err := aq.Update(ctx, updateInput1)
	if err != nil {
		t.Fatalf("First update failed: %v", err)
	}
	if updatedArticle.Title != newTitle {
		t.Errorf("Expected title %s after first update, got %s", newTitle, updatedArticle.Title)
	}

	newLikes := 20
	newReposts := 5
	updateInput2 := ArticleUpdateInput{
		Likes:     &newLikes,
		Reposts:   &newReposts,
		UpdatedAt: time.Now().UTC(),
	}
	aq.filters = bson.M{}
	aq.FilterID(articleID)
	updatedArticle, err = aq.Update(ctx, updateInput2)
	if err != nil {
		t.Fatalf("Second update failed: %v", err)
	}
	if updatedArticle.Likes != newLikes {
		t.Errorf("Expected likes %d after second update, got %d", newLikes, updatedArticle.Likes)
	}
	if updatedArticle.Reposts != newReposts {
		t.Errorf("Expected reposts %d after second update, got %d", newReposts, updatedArticle.Reposts)
	}
}

//AUTHORS TESTS

// setupAuthors создаёт экземпляр AuthorsQ для тестовой базы и очищает коллекцию.
func setupAuthors(t *testing.T) *AuthorsQ {
	aq, err := NewAuthors("test", "mongodb://localhost:7100")
	if err != nil {
		t.Fatalf("failed to create AuthorsQ: %v", err)
	}
	if err := aq.collection.Drop(context.Background()); err != nil {
		t.Fatalf("failed to drop test collection: %v", err)
	}
	aq.filters = bson.M{}
	aq.limit = 0
	aq.skip = 0
	aq.sort = bson.D{}
	return aq
}

func TestAuthors_Insert_Get_Count_Delete(t *testing.T) {
	aq := setupAuthors(t)
	ctx := context.Background()

	authorID := uuid.New()
	now := time.Now().UTC()
	input := AuthorInsertInput{
		ID:        authorID,
		Name:      "John Doe",
		Desc:      nil,
		Avatar:    nil,
		Email:     nil,
		Telegram:  nil,
		Twitter:   nil,
		CreatedAt: now,
	}

	if err := aq.Insert(ctx, input); err != nil {
		t.Fatalf("failed to insert author: %v", err)
	}

	count, err := aq.Count(ctx)
	if err != nil {
		t.Fatalf("failed to count authors: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	aq.New().FilterID(authorID)
	author, err := aq.Get(ctx)
	if err != nil {
		t.Fatalf("failed to get author: %v", err)
	}
	if author.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", author.Name)
	}

	if err := aq.FilterID(authorID).Delete(ctx); err != nil {
		t.Fatalf("failed to delete author: %v", err)
	}
	count, err = aq.Count(ctx)
	if err != nil {
		t.Fatalf("failed to count authors after deletion: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0 after delete, got %d", count)
	}
}

func TestAuthors_FilterName_Sort_Limit_Skip(t *testing.T) {
	aq := setupAuthors(t)
	ctx := context.Background()

	authors := []AuthorInsertInput{
		{
			ID:        uuid.New(),
			Name:      "Alice Smith",
			CreatedAt: time.Now().Add(-10 * time.Minute),
		},
		{
			ID:        uuid.New(),
			Name:      "Bob Doe",
			CreatedAt: time.Now().Add(-5 * time.Minute),
		},
		{
			ID:        uuid.New(),
			Name:      "Alice Doe",
			CreatedAt: time.Now(),
		},
	}

	for i, input := range authors {
		if input.Name == "" {
			input.Name = fmt.Sprintf("Author %d", i)
		}
		input.CreatedAt = time.Now().Add(time.Duration(i) * time.Minute)
		if err := aq.Insert(ctx, input); err != nil {
			t.Fatalf("failed to insert author %d: %v", i, err)
		}
	}

	aq = aq.New()
	aq.FilterName("Doe")
	results, err := aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with FilterName failed: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 authors with 'Doe' in name, got %d", len(results))
	}

	aq = aq.New()
	aq.Sort("created_at", true)
	results, err = aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with Sort failed: %v", err)
	}
	if len(results) < 3 {
		t.Errorf("expected at least 3 authors, got %d", len(results))
	} else {
		if results[0].CreatedAt.After(results[1].CreatedAt) {
			t.Errorf("authors not sorted ascending by CreatedAt")
		}
	}

	aq = aq.New()
	aq.Sort("created_at", true)
	aq.Limit(1)
	aq.Skip(1)
	results, err = aq.Select(ctx)
	if err != nil {
		t.Fatalf("Select with Limit/Skip failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 author after applying limit and skip, got %d", len(results))
	}
}

func TestAuthors_Update(t *testing.T) {
	aq := setupAuthors(t)
	ctx := context.Background()

	authorID := uuid.New()
	now := time.Now().UTC()
	input := AuthorInsertInput{
		ID:        authorID,
		Name:      "Original Name",
		CreatedAt: now,
	}
	if err := aq.New().Insert(ctx, input); err != nil {
		t.Fatalf("failed to insert author: %v", err)
	}

	newName := "Updated Name"
	newDesc := "Updated description"
	updateInput := AuthorUpdateInput{
		Name:      &newName,
		Desc:      &newDesc,
		UpdatedAt: time.Now().UTC(),
	}
	updatedAuthor, err := aq.New().FilterID(authorID).Update(ctx, updateInput)
	if err != nil {
		t.Fatalf("failed to update author: %v", err)
	}
	if updatedAuthor.Name != newName {
		t.Errorf("expected name '%s', got '%s'", newName, updatedAuthor.Name)
	}
	if updatedAuthor.Desc == nil || *updatedAuthor.Desc != newDesc {
		t.Errorf("expected desc '%s', got '%v'", newDesc, updatedAuthor.Desc)
	}
}
