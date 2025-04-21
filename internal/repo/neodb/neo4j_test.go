package neodb

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	testURI  = "bolt://localhost:7687"
	testUser = "neo4j"
	testPass = "password"
)

// cleanDB drops all nodes and relationships.
// It runs before each test to have a fresh start.
func cleanDB(t *testing.T, drv neo4j.Driver) {
	t.Helper()
	sess, err := drv.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		t.Fatalf("open session: %v", err)
	}
	defer sess.Close()

	_, err = sess.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		// remove everything
		_, err := tx.Run("MATCH (n) DETACH DELETE n", nil)
		return nil, err
	})
	if err != nil {
		t.Fatalf("clean db: %v", err)
	}
}

// TestArticleCRUD covers Create, GetByID, UpdateStatus, Delete.
func TestArticleCRUD(t *testing.T) {
	repo, err := NewArticles(testURI, testUser, testPass)
	if err != nil {
		t.Fatalf("NewArticles fail: %v", err)
	}
	cleanDB(t, repo.driver)

	ctx := context.Background()
	id := uuid.New()

	// 1. Create node
	if _, err := repo.Create(ctx, ArticleInsertInput{ID: id, Status: enums.ArticleStatusActive}); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// 2. Read it
	a, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if a.ID != id || a.Status != enums.ArticleStatusActive {
		t.Errorf("GetByID wrong data: %+v", a)
	}

	// 3. Update status using new Update()
	newStatus := enums.ArticleStatusInactive
	updated, err := repo.Update(ctx, id, ArticleUpdateInput{Status: &newStatus})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.ID != id {
		t.Errorf("Updated.ID changed: got %v want %v", updated.ID, id)
	}
	if updated.Status != newStatus {
		t.Errorf("Status not updated: got %v want %v", updated.Status, newStatus)
	}

	// 4. Delete node
	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := repo.GetByID(ctx, id); err == nil {
		t.Errorf("Expected error after delete, got nil")
	}
}

// TestAuthorCRUD covers Create, GetByID, UpdateStatus, Delete.
func TestAuthorCRUD(t *testing.T) {
	repo, err := NewAuthors(testURI, testUser, testPass)
	if err != nil {
		t.Fatalf("NewAuthors fail: %v", err)
	}
	cleanDB(t, repo.driver)

	ctx := context.Background()
	id := uuid.New()

	// 1. Create node
	if _, err := repo.Create(ctx, AuthorCreateInput{ID: id, Status: enums.AuthorStatusActive}); err != nil {
		t.Fatalf("Author Create failed: %v", err)
	}

	// 2. Read it
	au, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Author GetByID failed: %v", err)
	}
	if au.ID != id || au.Status != enums.AuthorStatusActive {
		t.Errorf("GetByID wrong: %+v", au)
	}

	// 3. Update status using new Update()
	newStatus := enums.AuthorStatusInactive
	updated, err := repo.Update(ctx, id, AuthorUpdateInput{Status: &newStatus})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.ID != id {
		t.Errorf("Updated.ID changed: got %v want %v", updated.ID, id)
	}
	if updated.Status != newStatus {
		t.Errorf("Status not updated: got %v want %v", updated.Status, newStatus)
	}

	// 4. Delete node
	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("Author Delete failed: %v", err)
	}
	if _, err := repo.GetByID(ctx, id); err == nil {
		t.Errorf("Expected error after delete, got nil")
	}
}

// TestAuthorshipRelation covers creating and deleting relation.
func TestAuthorshipRelation(t *testing.T) {
	// prepare article and author
	artRepo, _ := NewArticles(testURI, testUser, testPass)
	authRepo, _ := NewAuthors(testURI, testUser, testPass)
	asRepo, _ := NewAuthorship(testURI, testUser, testPass)
	cleanDB(t, artRepo.driver)

	ctx := context.Background()
	aID := uuid.New()
	uID := uuid.New()

	// create nodes
	if _, err := artRepo.Create(ctx, ArticleInsertInput{ID: aID, Status: enums.ArticleStatusActive}); err != nil {
		t.Fatalf("create article: %v", err)
	}
	if _, err := authRepo.Create(ctx, AuthorCreateInput{ID: uID, Status: enums.AuthorStatusActive}); err != nil {
		t.Fatalf("create author: %v", err)
	}

	// create relationship
	if err := asRepo.Create(ctx, aID, uID); err != nil {
		t.Fatalf("Authorship Create failed: %v", err)
	}

	// get for article
	authors, err := asRepo.GetForArticle(ctx, aID)
	if err != nil {
		t.Fatalf("GetForArticle failed: %v", err)
	}
	if len(authors) != 1 || authors[0] != uID {
		t.Errorf("GetForArticle wrong: %v", authors)
	}

	// get for author
	arts, err := asRepo.GetForAuthor(ctx, uID)
	if err != nil {
		t.Fatalf("GetForAuthor failed: %v", err)
	}
	if len(arts) != 1 || arts[0] != aID {
		t.Errorf("GetForAuthor wrong: %v", arts)
	}

	// delete relation
	if err := asRepo.Delete(ctx, aID, uID); err != nil {
		t.Errorf("Authorship Delete failed: %v", err)
	}
	after, _ := asRepo.GetForArticle(ctx, aID)
	if len(after) != 0 {
		t.Errorf("Relation not removed: %v", after)
	}
}

// TestHashtagRelation covers Create, GetForArticle, GetArticlesForTag, Delete.
func TestHashtagRelation(t *testing.T) {
	artRepo, _ := NewArticles(testURI, testUser, testPass)
	tagRepo, _ := NewHashtag(testURI, testUser, testPass)
	// we need also tag nodes, use TagsImpl
	tagsImpl, _ := NewTags(testURI, testUser, testPass)

	cleanDB(t, artRepo.driver)

	ctx := context.Background()
	aID := uuid.New()
	tagName := "golang"

	// create article and tag nodes
	_, err := artRepo.Create(ctx, ArticleInsertInput{ID: aID, Status: enums.ArticleStatusActive})
	if err != nil {
		t.Fatalf("Article Create failed: %v", err)
	}

	_, err = tagsImpl.Create(ctx, TagCreateInput{
		Name:      tagName,
		Status:    enums.TagStatusActive,
		Type:      enums.TagTypeTopic,
		Color:     "blue",
		Icon:      "icon.png",
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}

	// add hashtag
	if err := tagRepo.Create(ctx, aID, tagName); err != nil {
		t.Fatalf("Hashtag Create failed: %v", err)
	}

	// get for article
	hts, err := tagRepo.GetForArticle(ctx, aID)
	if err != nil {
		t.Fatalf("GetForArticle tag failed: %v", err)
	}
	if len(hts) != 1 || hts[0] != tagName {
		t.Errorf("Unexpected tags: %v", hts)
	}

	// get articles for tag
	arts, err := tagRepo.GetArticlesForTag(ctx, tagName)
	if err != nil {
		t.Fatalf("GetArticlesForTag failed: %v", err)
	}
	if len(arts) != 1 || arts[0] != aID {
		t.Errorf("Unexpected articles: %v", arts)
	}

	// delete hashtag
	if err := tagRepo.Delete(ctx, aID, tagName); err != nil {
		t.Errorf("Hashtag Delete failed: %v", err)
	}
	empty, _ := tagRepo.GetForArticle(ctx, aID)
	if len(empty) != 0 {
		t.Errorf("Hashtag not removed: %v", empty)
	}
}

// TestSetForArticle tests bulk replace of authors and tags.
func TestSetForArticle(t *testing.T) {
	artRepo, _ := NewArticles(testURI, testUser, testPass)
	authImpl, _ := NewAuthorship(testURI, testUser, testPass)
	tagRepo, _ := NewHashtag(testURI, testUser, testPass)
	notTagsImpl, _ := NewTags(testURI, testUser, testPass)
	cleanDB(t, artRepo.driver)

	ctx := context.Background()
	aID := uuid.New()
	// create article, authors, tags
	_, _ = artRepo.Create(ctx, ArticleInsertInput{ID: aID, Status: enums.ArticleStatusActive})
	authors := []uuid.UUID{uuid.New(), uuid.New()}
	for _, u := range authors {
		authorsQ, err := NewAuthors(testURI, testUser, testPass)
		if err != nil {
			t.Fatalf("NewAuthors failed: %v", err)
		}

		if _, err := authorsQ.Create(ctx, AuthorCreateInput{
			ID:     u,
			Status: enums.AuthorStatusActive,
		}); err != nil {
			t.Fatalf("Create authors failed: %v", err)
		}
	}
	tags := []string{"t1", "t2"}
	for _, tg := range tags {
		notTagsImpl.Create(ctx, TagCreateInput{
			Name: tg, Status: enums.TagStatusActive,
			Type: enums.TagTypeDefault, Color: "c", Icon: "i", CreatedAt: time.Now().UTC(),
		})
	}

	// set authors
	if err := authImpl.SetForArticle(ctx, aID, authors); err != nil {
		t.Fatalf("SetForArticle authors failed: %v", err)
	}
	gotAuth, _ := authImpl.GetForArticle(ctx, aID)
	if len(gotAuth) != 2 {
		t.Errorf("Expected 2 authors, got %v", gotAuth)
	}

	// set tags
	if err := tagRepo.SetForArticle(ctx, aID, tags); err != nil {
		t.Fatalf("SetForArticle tags failed: %v", err)
	}
	gotTags, _ := tagRepo.GetForArticle(ctx, aID)
	if len(gotTags) != 2 {
		t.Errorf("Expected 2 tags, got %v", gotTags)
	}
}

// TestTagsCRUD covers Create, Get, GetAll, Update, Delete.
func TestTagsCRUD(t *testing.T) {
	timpl, err := NewTags(testURI, testUser, testPass)
	if err != nil {
		t.Fatalf("NewTags failed: %v", err)
	}
	cleanDB(t, timpl.driver)

	ctx := context.Background()
	name := "tagx"
	now := time.Now().UTC()

	// create
	_, err = timpl.Create(ctx, TagCreateInput{
		Name: name, Status: enums.TagStatusActive,
		Type: enums.TagTypeTopic, Color: "red",
		Icon: "icon", CreatedAt: now,
	})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}

	// get by name
	tag, err := timpl.Get(ctx, "tag")
	if err != nil {
		t.Fatalf("Tag Get failed: %v", err)
	}
	if tag.Name != name {
		t.Errorf("Got wrong tag: %s", tag.Name)
	}

	// get all (popularity sort); just one
	all, err := timpl.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(all))
	}

	// update fields
	newIcon := "newicon"
	updatedAt := time.Now().Add(time.Minute).UTC()
	updated, err := timpl.Update(ctx, name, TagUpdateInput{
		Icon:      &newIcon,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		t.Fatalf("Tag Update failed: %v", err)
	}
	if updated.Icon != newIcon {
		t.Errorf("Icon not updated: %s", updated.Icon)
	}

	// delete
	if err := timpl.Delete(ctx, name); err != nil {
		t.Errorf("Tag Delete failed: %v", err)
	}
	if _, err := timpl.Get(ctx, name); err == nil {
		t.Errorf("Expected error after delete")
	}
}

func TestHashtagEmpty(t *testing.T) {
	artRepo, _ := NewArticles(testURI, testUser, testPass)
	tagRepo, _ := NewHashtag(testURI, testUser, testPass)
	cleanDB(t, artRepo.driver)

	aID := uuid.New()
	got, err := tagRepo.GetForArticle(context.Background(), aID)
	if err != nil {
		t.Fatalf("GetForArticle error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Expected no tags, got %v", got)
	}

	got2, err := tagRepo.GetArticlesForTag(context.Background(), "none")
	if err != nil {
		t.Fatalf("GetArticlesForTag error: %v", err)
	}
	if len(got2) != 0 {
		t.Errorf("Expected no articles, got %v", got2)
	}
}

// TestHashtagSetEmpty tests SetForArticle with empty slice deletes all tags.
func TestHashtagSetEmpty(t *testing.T) {
	artRepo, _ := NewArticles(testURI, testUser, testPass)
	tagsImpl, _ := NewTags(testURI, testUser, testPass)
	tagRepo, _ := NewHashtag(testURI, testUser, testPass)
	cleanDB(t, artRepo.driver)

	ctx := context.Background()
	aID := uuid.New()
	_, err := artRepo.Create(ctx, ArticleInsertInput{ID: aID, Status: enums.ArticleStatusActive})
	// create two tags
	tags := []string{"t1", "t2"}
	for _, tg := range tags {
		_, err = tagsImpl.Create(ctx, TagCreateInput{
			Name: tg, Status: enums.TagStatusActive,
			Type: enums.TagTypeDefault, Color: "c", Icon: "i", CreatedAt: time.Now().UTC(),
		})
	}
	// set tags
	_ = tagRepo.SetForArticle(ctx, aID, tags)
	if got, _ := tagRepo.GetForArticle(ctx, aID); len(got) != 2 {
		t.Fatalf("Expected 2 tags, got %v", got)
	}
	// clear tags
	_ = tagRepo.SetForArticle(ctx, aID, []string{})
	got2, err := tagRepo.GetForArticle(ctx, aID)
	if err != nil {
		t.Fatalf("GetForArticle error: %v", err)
	}
	if len(got2) != 0 {
		t.Errorf("Expected no tags after clear, got %v", got2)
	}
}

// TestTagsPopularitySort tests GetAll orders by relationship count desc.
func TestTagsPopularitySort(t *testing.T) {
	artRepo, _ := NewArticles(testURI, testUser, testPass)
	tagsImpl, _ := NewTags(testURI, testUser, testPass)
	tagRepo, _ := NewHashtag(testURI, testUser, testPass)
	cleanDB(t, artRepo.driver)

	ctx := context.Background()
	// create article and 2 tags
	a1, a2 := uuid.New(), uuid.New()
	_, err := artRepo.Create(ctx, ArticleInsertInput{ID: a1, Status: enums.ArticleStatusActive})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}
	_, err = artRepo.Create(ctx, ArticleInsertInput{ID: a2, Status: enums.ArticleStatusActive})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}
	t1, t2 := "tagA", "tagB"
	_, err = tagsImpl.Create(ctx, TagCreateInput{Name: t1, Status: enums.TagStatusActive, Type: enums.TagTypeDefault, Color: "c", Icon: "i", CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}
	_, err = tagsImpl.Create(ctx, TagCreateInput{Name: t2, Status: enums.TagStatusActive, Type: enums.TagTypeDefault, Color: "c", Icon: "i", CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}
	// assign t1 to both articles, t2 to one
	_ = tagRepo.Create(ctx, a1, t1)
	_ = tagRepo.Create(ctx, a2, t1)
	_ = tagRepo.Create(ctx, a1, t2)
	// now GetAll should return [t1, t2]
	all, err := tagsImpl.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll error: %v", err)
	}
	if len(all) < 2 || all[0].Name != t1 || all[1].Name != t2 {
		t.Errorf("Expected order [%s,%s], got %v", t1, t2, []string{all[0].Name, all[1].Name})
	}
}

// TestTagsGetContains tests Get is case-insensitive contains.
func TestTagsGetContains(t *testing.T) {
	timpl, _ := NewTags(testURI, testUser, testPass)
	cleanDB(t, timpl.driver)

	ctx := context.Background()
	_, err := timpl.Create(ctx, TagCreateInput{Name: "GoLang", Status: enums.TagStatusActive, Type: enums.TagTypeDefault, Color: "c", Icon: "i", CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}
	// search lowercase fragment
	got, err := timpl.Get(ctx, "lang")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got.Name != "GoLang" {
		t.Errorf("Expected GoLang, got %s", got.Name)
	}
}

// TestTagUpdateMultipleFields tests updating all fields at once.
func TestTagUpdateMultipleFields(t *testing.T) {
	timpl, _ := NewTags(testURI, testUser, testPass)
	cleanDB(t, timpl.driver)

	ctx := context.Background()
	name := "old"
	now := time.Now().UTC()
	_, err := timpl.Create(ctx, TagCreateInput{
		Name: name, Status: enums.TagStatusActive,
		Type: enums.TagTypeDefault, Color: "red",
		Icon: "ico", CreatedAt: now,
	})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}

	newName := "new"
	newStatus := enums.TagStatusInactive
	newType := enums.TagTypeTopic
	newColor := "blue"
	newIcon := "ico2"
	updatedAt := now.Add(time.Minute)

	updated, err := timpl.Update(ctx, name, TagUpdateInput{
		NewName:   &newName,
		Status:    &newStatus,
		Type:      &newType,
		Color:     &newColor,
		Icon:      &newIcon,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if updated.Name != newName ||
		updated.Status != newStatus ||
		updated.Type != newType ||
		updated.Color != newColor ||
		updated.Icon != newIcon ||
		updated.UpdatedAt == nil || !updated.UpdatedAt.Equal(updatedAt) {
		t.Errorf("Fields not all updated: %+v", updated)
	}
}

// TestTagUpdateNoop tests Update with no change returns existing.
func TestTagUpdateNoop(t *testing.T) {
	timpl, err := NewTags(testURI, testUser, testPass)
	if err != nil {
		t.Fatalf("NewTags failed: %v", err)
	}

	ctx := context.Background()
	name := "keep"
	now := time.Now().UTC()
	cr, err := timpl.Create(ctx, TagCreateInput{Name: name, Status: enums.TagStatusActive, Type: enums.TagTypeDefault, Color: "x", Icon: "y", CreatedAt: now})
	if err != nil {
		t.Fatalf("Tag Create failed: %v", err)
	}

	newTime := now.Add(2 * time.Minute)
	up, err := timpl.Update(ctx, name, TagUpdateInput{UpdatedAt: newTime})
	if err != nil {
		t.Fatalf("Update noop error: %v", err)
	}

	if up.Name != cr.Name {
		t.Errorf("Fields changed: %+v", up)
	}
}
