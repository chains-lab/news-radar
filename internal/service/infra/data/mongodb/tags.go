package mongodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/service/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tags interface {
	New() Tags

	Insert(ctx context.Context, tag *models.Tag) (*models.Tag, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]*models.Tag, error)
	Get(ctx context.Context) (*models.Tag, error)

	FiltersID(id uuid.UUID) Tags
	FiltersName(name string) Tags

	Update(ctx context.Context, fields map[string]any) (*models.Tag, error)

	Limit(limit int64) Tags
	Skip(skip int64) Tags
	Sort(field string, ascending bool) Tags
}

type tags struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection

	filters bson.M
	sort    bson.D
	limit   int64
	skip    int64
}

func NewTags(uri, dbName, collectionName string) (Tags, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	database := client.Database(dbName)
	coll := database.Collection(collectionName)

	return &tags{
		client:     client,
		database:   database,
		collection: coll,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}, nil
}

func (t *tags) New() Tags {
	return &tags{
		client:     t.client,
		database:   t.database,
		collection: t.collection,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}
}

func (t *tags) Insert(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	_, err := t.collection.InsertOne(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tag: %w", err)
	}
	return tag, nil
}

func (t *tags) Delete(ctx context.Context) error {
	_, err := t.collection.DeleteOne(ctx, t.filters)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

func (t *tags) Count(ctx context.Context) (int64, error) {
	return t.collection.CountDocuments(ctx, t.filters)
}

func (t *tags) Select(ctx context.Context) ([]*models.Tag, error) {
	cursor, err := t.collection.Find(ctx, t.filters)
	if err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}
	defer cursor.Close(ctx)

	var tagsList []*models.Tag
	if err = cursor.All(ctx, &tagsList); err != nil {
		return nil, fmt.Errorf("failed to decode tags: %w", err)
	}
	return tagsList, nil
}

func (t *tags) Get(ctx context.Context) (*models.Tag, error) {
	var tag models.Tag
	err := t.collection.FindOne(ctx, t.filters).Decode(&tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

func (t *tags) FiltersID(id uuid.UUID) Tags {
	t.filters["_id"] = id
	return t
}

func (t *tags) FiltersName(name string) Tags {
	t.filters["name"] = bson.M{
		"$regex":   fmt.Sprintf(".*%s.*", name),
		"$options": "i",
	}
	return t
}

func (t *tags) Update(ctx context.Context, fields map[string]any) (*models.Tag, error) {
	validFields := map[string]bool{
		"name":     true,
		"status":   true,
		"category": true,
	}
	updateFields := bson.M{}
	for key, value := range fields {
		if validFields[key] {
			updateFields[key] = value
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Tag
	err := t.collection.FindOneAndUpdate(ctx, t.filters, bson.M{"$set": updateFields}, opts).Decode(&updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	for key, value := range updateFields {
		if _, exists := t.filters[key]; exists {
			t.filters[key] = value
		}
	}

	return &updated, nil
}

func (t *tags) Limit(limit int64) Tags {
	t.limit = limit
	return t
}

func (t *tags) Skip(skip int64) Tags {
	t.skip = skip
	return t
}

func (t *tags) Sort(field string, ascending bool) Tags {
	order := 1
	if !ascending {
		order = -1
	}
	t.sort = bson.D{{field, order}}
	return t
}
