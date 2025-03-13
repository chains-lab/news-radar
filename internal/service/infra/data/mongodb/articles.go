package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Articles interface {
	New() Articles

	Insert(ctx context.Context, article *models.Article) (*models.Article, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]*models.Article, error)
	Get(ctx context.Context) (*models.Article, error)

	FilterID(id uuid.UUID) Articles
	FilterTitle(name string) Articles
	FilterAuthors(authors ...uuid.UUID) Articles
	FilterTags(tags ...string) Articles
	FilterDate(filters map[string]any, after bool) Articles

	Update(ctx context.Context, fields map[string]any) (*models.Article, error)

	Limit(limit int64) Articles
	Skip(skip int64) Articles
	Sort(field string, ascending bool) Articles
}

type articles struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection

	filters bson.M
	sort    bson.D
	limit   int64
	skip    int64
}

func NewArticles(uri, dbName, collectionName string) (Articles, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	database := client.Database(dbName)
	coll := database.Collection(collectionName)

	return &articles{
		client:     client,
		database:   database,
		collection: coll,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}, nil
}

func (a *articles) New() Articles {
	return &articles{
		client:     a.client,
		database:   a.database,
		collection: a.collection,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}
}

func (a *articles) Insert(ctx context.Context, article *models.Article) (*models.Article, error) {
	_, err := a.collection.InsertOne(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("failed to insert article: %w", err)
	}
	return article, nil
}

func (a *articles) Delete(ctx context.Context) error {
	_, err := a.collection.DeleteOne(ctx, a.filters)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}
	return nil
}

func (a *articles) Count(ctx context.Context) (int64, error) {
	return a.collection.CountDocuments(ctx, a.filters)
}

func (a *articles) Select(ctx context.Context) ([]*models.Article, error) {
	findOptions := options.Find()
	if a.limit > 0 {
		findOptions.SetLimit(a.limit)
	}
	if a.skip > 0 {
		findOptions.SetSkip(a.skip)
	}
	if len(a.sort) > 0 {
		findOptions.SetSort(a.sort)
	}

	cursor, err := a.collection.Find(ctx, a.filters, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to select articles: %w", err)
	}
	defer cursor.Close(ctx)

	var arts []*models.Article
	if err = cursor.All(ctx, &arts); err != nil {
		return nil, fmt.Errorf("failed to decode articles: %w", err)
	}
	return arts, nil
}

func (a *articles) Get(ctx context.Context) (*models.Article, error) {
	var art models.Article
	err := a.collection.FindOne(ctx, a.filters).Decode(&art)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return &art, nil
}

func (a *articles) FilterID(id uuid.UUID) Articles {
	a.filters["_id"] = id
	return a
}

func (a *articles) FilterTitle(title string) Articles {
	a.filters["title"] = bson.M{
		"$regex":   fmt.Sprintf(".*%s.*", title),
		"$options": "i",
	}
	return a
}

func (a *articles) FilterAuthors(authors ...uuid.UUID) Articles {
	if len(authors) == 0 {
		return a
	}
	a.filters["authors"] = bson.M{
		"$in": authors,
	}
	return a
}

func (a *articles) FilterTags(tags ...string) Articles {
	if len(tags) == 0 {
		return a
	}
	a.filters["tags"] = bson.M{
		"$in": tags,
	}
	return a
}

func (a *articles) FilterDate(filters map[string]any, after bool) Articles {
	validDateFields := map[string]bool{
		"updated_at": true,
		"closed_at":  true,
	}

	var op string
	if after {
		op = "$gte"
	} else {
		op = "$lte"
	}

	for field, value := range filters {
		if !validDateFields[field] {
			continue
		}
		if value == nil {
			continue
		}

		var t time.Time
		switch val := value.(type) {
		case time.Time:
			t = val
		case *time.Time:
			t = *val
		case string:
			parsed, err := time.Parse(time.RFC3339, val)
			if err != nil {
				continue
			}
			t = parsed
		default:
			continue
		}

		a.filters[field] = bson.M{op: t}
	}

	return a
}

func (a *articles) Update(ctx context.Context, fields map[string]any) (*models.Article, error) {
	validFields := map[string]bool{
		"title":       true,
		"icon":        true,
		"description": true,
		"authors":     true,
		"content":     true,
		"likes":       true,
		"reposts":     true,
		"tags":        true,
		"updated_at":  true,
	}
	updateFields := bson.M{}
	for key, value := range fields {
		if validFields[key] {
			updateFields[key] = value
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Article
	err := a.collection.FindOneAndUpdate(ctx, a.filters, bson.M{"$set": updateFields}, opts).Decode(&updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}
	return &updated, nil
}

func (a *articles) Limit(limit int64) Articles {
	a.limit = limit
	return a
}

func (a *articles) Skip(skip int64) Articles {
	a.skip = skip
	return a
}

func (a *articles) Sort(field string, ascending bool) Articles {
	order := 1
	if !ascending {
		order = -1
	}
	a.sort = bson.D{{Key: field, Value: order}}
	return a
}
