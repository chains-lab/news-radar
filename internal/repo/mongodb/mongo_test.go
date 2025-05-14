package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chains-lab/news-radar/internal/content"
	"github.com/chains-lab/news-radar/internal/enums"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// setupArticle и setupAuthors находятся в этом же пакете
func setupArticle(t *testing.T) *ArticlesQ {
	aq, err := NewArticles("test", "mongodb://localhost:7100")
	if err != nil {
		t.Fatalf("failed to create ArticlesQ: %v", err)
	}
	// очистка тестовой коллекции
	if err := aq.collection.Drop(context.Background()); err != nil {
		t.Fatalf("failed to drop test collection: %v", err)
	}
	// сброс параметров запроса
	aq.filters = bson.M{}
	aq.limit = 0
	aq.skip = 0
	aq.sort = bson.D{}
	return aq
}

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

func TestArticlesFilterTitle(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	// insert two articles with different titles
	id1 := uuid.New()
	id2 := uuid.New()
	now := time.Now().UTC()
	_ = aq.Insert(ctx, ArticleInsertInput{ID: id1, Title: "Go tutorial", CreatedAt: now})
	_ = aq.Insert(ctx, ArticleInsertInput{ID: id2, Title: "Python guide", CreatedAt: now})

	// filter for "Go"
	list, err := aq.New().FilterTitle("Go").Select(ctx)
	if err != nil {
		t.Fatalf("FilterTitle failed: %v", err)
	}
	if len(list) != 1 || list[0].ID != id1 {
		t.Errorf("FilterTitle returned wrong result: %+v", list)
	}
}

func TestArticlesPagination(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	// insert 5 articles
	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		id := uuid.New()
		title := fmt.Sprintf("Art%d", i)
		_ = aq.Insert(ctx, ArticleInsertInput{ID: id, Title: title, CreatedAt: now})
	}

	// get page 2 with 2 items per page
	page, err := aq.New().Sort("title", true).Limit(2).Skip(2).Select(ctx)
	if err != nil {
		t.Fatalf("Pagination Select failed: %v", err)
	}
	if len(page) != 2 {
		t.Errorf("Expected 2 items on page, got %d", len(page))
	}
}

func TestArticlesFilterDate(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	id := uuid.New()
	created := time.Now().Add(-10 * time.Minute).UTC()
	_ = aq.Insert(ctx, ArticleInsertInput{ID: id, Title: "Old", CreatedAt: created})

	// update article to set UpdatedAt to now
	now := time.Now().UTC()
	_, _ = aq.New().FilterID(id).Update(ctx, ArticleUpdateInput{Title: ptrString("New"), UpdatedAt: now})

	// find articles updated after 5 minutes ago
	filters := map[string]any{"updated_at": now.Add(-5 * time.Minute)}
	list, err := aq.New().FilterDate(filters, true).Select(ctx)
	if err != nil {
		t.Fatalf("FilterDate failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 updated article, got %d", len(list))
	}
}

func TestArticlesFilterStatus(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	id := uuid.New()
	_ = aq.Insert(ctx, ArticleInsertInput{ID: id, Title: "Stat", CreatedAt: time.Now().UTC()})
	// status is always Active on insert
	list, err := aq.New().FilterStatus(enums.ArticleStatusPublished).Select(ctx)
	if err != nil {
		t.Fatalf("FilterStatus failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 active article, got %d", len(list))
	}
}

func TestArticlesBulkSelect(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	ids := make([]uuid.UUID, 3)
	now := time.Now().UTC()
	for i := range ids {
		ids[i] = uuid.New()
		_ = aq.Insert(ctx, ArticleInsertInput{ID: ids[i], Title: fmt.Sprintf("A%d", i), CreatedAt: now})
	}
	all, err := aq.New().Select(ctx)
	if err != nil {
		t.Fatalf("Bulk Select failed: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Expected 3 articles, got %d", len(all))
	}
}

func TestArticlesCrud(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	// Insert
	id := uuid.New()
	title := "Test Article"
	created := time.Now().UTC()
	if err := aq.Insert(ctx, ArticleInsertInput{ID: id, Title: title, CreatedAt: created}); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Count
	count, err := aq.Count(ctx)
	if err != nil || count != 1 {
		t.Fatalf("Count expected 1, got %d, err %v", count, err)
	}

	// Get
	aqt := aq.New().FilterID(id)
	art, err := aqt.Get(ctx)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if art.ID != id || art.Title != title {
		t.Errorf("Get returned wrong article: %v", art)
	}

	// Update Name
	newTitle := "Updated"
	updated := time.Now().UTC()
	tstArt, err := aqt.Update(ctx, ArticleUpdateInput{Title: &newTitle, UpdatedAt: updated})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if tstArt.Title != newTitle {
		t.Errorf("Update did not change title")
	}

	// Delete
	if err := aqt.Delete(ctx); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	count, _ = aq.Count(ctx)
	if count != 0 {
		t.Errorf("Delete did not remove article, count=%d", count)
	}
}

func TestArticlesSort(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()

	// insert titles out of order
	_ = aq.Insert(ctx, ArticleInsertInput{ID: uuid.New(), Title: "B", CreatedAt: time.Now().UTC()})
	_ = aq.Insert(ctx, ArticleInsertInput{ID: uuid.New(), Title: "A", CreatedAt: time.Now().UTC()})

	asc, _ := aq.New().Sort("title", true).Select(ctx)
	if len(asc) >= 2 && asc[0].Title != "A" {
		t.Errorf("Expected first title A, got %s", asc[0].Title)
	}
	desc, _ := aq.New().Sort("title", false).Select(ctx)
	if len(desc) >= 2 && desc[0].Title != "B" {
		t.Errorf("Expected first title B, got %s", desc[0].Title)
	}
}

// ptrString помогает строить текстовые блоки
func ptrString(s string) *string { return &s }

func TestUpdateContent(t *testing.T) {
	aq := setupArticle(t)
	ctx := context.Background()
	id := uuid.New()
	created := time.Now().UTC()
	// вставляем начальный документ
	if err := aq.Insert(ctx, ArticleInsertInput{ID: id, Title: "C", CreatedAt: created}); err != nil {
		t.Fatalf("Insert for UpdateContentSection failed: %v", err)
	}
	aqt := aq.New().FilterID(id)

	// 1. Добавление текстовой секции
	tb := content.TextBlock{Text: ptrString("Hello"), Marks: nil, Color: nil, Link: nil}
	secText := content.Section{ID: "sec1", Type: enums.SectionTypeText, Text: []content.TextBlock{tb}}
	art1, err := aqt.UpdateContentSection(ctx, 0, secText, time.Now().UTC())
	if err != nil {
		t.Fatalf("Add text section failed: %v", err)
	}
	if len(art1.Content) != 1 || art1.Content[0].ID != "sec1" || len(art1.Content[0].Text) != 1 || *art1.Content[0].Text[0].Text != "Hello" {
		t.Errorf("Unexpected content after add: %+v", art1.Content)
	}

	// 2. Обновление на медиа секцию
	media := &content.Media{URL: "u", Caption: "c", Alt: "a", Width: 1, Height: 2, Source: "s"}
	secMedia := content.Section{ID: "sec1", Type: enums.SectionTypeMedia, Media: media}
	art2, err := aqt.UpdateContentSection(ctx, 0, secMedia, time.Now().UTC())
	if err != nil {
		t.Fatalf("Update media section failed: %v", err)
	}
	if art2.Content[0].Media == nil || art2.Content[0].Media.URL != "u" {
		t.Errorf("Unexpected media after update: %+v", art2.Content[0].Media)
	}

	// 3. Добавление аудио секции в конец
	audio := &content.Audio{URL: "au", Duration: 3, Caption: "cap", Icon: "ic"}
	secAudio := content.Section{ID: "sec2", Type: enums.SectionTypeAudio, Audio: audio}
	art3, err := aqt.UpdateContentSection(ctx, 1, secAudio, time.Now().UTC())
	if err != nil {
		t.Fatalf("Append audio section failed: %v", err)
	}
	if len(art3.Content) != 2 || art3.Content[1].Audio == nil || art3.Content[1].Audio.URL != "au" {
		t.Errorf("Unexpected audio after append: %+v", art3.Content[1].Audio)
	}

	// 4. Удаление секции по индексу
	empty := content.Section{}
	art4, err := aqt.UpdateContentSection(ctx, 0, empty, time.Now().UTC())
	if err != nil {
		t.Fatalf("Remove section failed: %v", err)
	}
	if len(art4.Content) != 1 || art4.Content[0].ID != "sec2" {
		t.Errorf("Unexpected content after remove: %+v", art4.Content)
	}

	// 5. No-op для выхода за пределы: меняется только UpdatedAt
	prev := art4.UpdatedAt
	time.Sleep(1 * time.Millisecond)
	art5, err := aqt.UpdateContentSection(ctx, 5, empty, time.Now().UTC())
	if err != nil {
		t.Fatalf("No-op section failed: %v", err)
	}
	if art5.UpdatedAt == nil || !art5.UpdatedAt.After(*prev) {
		t.Errorf("Expected UpdatedAt to change, was %v", art5.UpdatedAt)
	}
}

func TestAuthorsCrud(t *testing.T) {
	aq := setupAuthors(t)
	ctx := context.Background()
	id := uuid.New()
	created := time.Now().UTC()
	name := "John"
	if err := aq.Insert(ctx, AuthorInsertInput{ID: id, Name: name, Status: enums.AuthorStatusActive, CreatedAt: created}); err != nil {
		t.Fatalf("Author Insert failed: %v", err)
	}

	// Count
	count, err := aq.Count(ctx)
	if err != nil || count != 1 {
		t.Fatalf("Expected count 1, got %d, err %v", count, err)
	}

	// Get
	aqt := aq.New().FilterID(id)
	auth, err := aqt.Get(ctx)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if auth.Name != name {
		t.Errorf("Get returned wrong name: %s", auth.Name)
	}

	// Update
	newName := "Jane"
	updated := time.Now().UTC()
	auth2, err := aqt.Update(ctx, AuthorUpdateInput{Name: &newName, UpdatedAt: updated})
	if err != nil {
		t.Fatalf("Author Update failed: %v", err)
	}
	if auth2.Name != newName {
		t.Errorf("Update did not change name")
	}

	// FilterName
	results, err := aq.New().FilterName("ane").Select(ctx)
	if err != nil || len(results) != 1 {
		t.Errorf("FilterName failed: %v, len=%d", err, len(results))
	}

	// Delete
	if err := aqt.Delete(ctx); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	count, _ = aq.Count(ctx)
	if count != 0 {
		t.Errorf("Delete did not remove author, count=%d", count)
	}
}

func TestAuthorsBulkInsert(t *testing.T) {
	aq := setupAuthors(t)
	ctx := context.Background()

	names := []string{"Ann", "Bob", "Cya"}
	for _, n := range names {
		_ = aq.Insert(ctx, AuthorInsertInput{ID: uuid.New(), Name: n, Status: enums.AuthorStatusActive, CreatedAt: time.Now().UTC()})
	}
	list, err := aq.New().Select(ctx)
	if err != nil {
		t.Fatalf("Authors bulk Select failed: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("Expected 3 authors, got %d", len(list))
	}
}
