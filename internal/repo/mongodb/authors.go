package mongodb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AuthorsCollection = "authors"
)

type AuthorsQ struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection

	filters bson.M
	sort    bson.D
	limit   int64
	skip    int64
}

func NewAuthors(uri, dbName string) (*AuthorsQ, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	database := client.Database(dbName)
	coll := database.Collection(AuthorsCollection)

	return &AuthorsQ{
		client:     client,
		database:   database,
		collection: coll,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}, nil
}

func (a *AuthorsQ) New() *AuthorsQ {
	return &AuthorsQ{
		client:     a.client,
		database:   a.database,
		collection: a.collection,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}
}

func (a *AuthorsQ) Insert(ctx context.Context, author modelsdb.AuthorMongo) (modelsdb.AuthorMongo, error) {
	_, err := a.collection.InsertOne(ctx, author)
	if err != nil {
		return modelsdb.AuthorMongo{}, err
	}

	return author, nil
}

func (a *AuthorsQ) Delete(ctx context.Context) error {
	_, err := a.collection.DeleteOne(ctx, a.filters)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthorsQ) Count(ctx context.Context) (int64, error) {
	return a.collection.CountDocuments(ctx, a.filters)
}

func (a *AuthorsQ) Select(ctx context.Context) ([]modelsdb.AuthorMongo, error) {
	cursor, err := a.collection.Find(ctx, a.filters)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var aths []modelsdb.AuthorMongo

	if err = cursor.All(ctx, &aths); err != nil {
		return nil, err
	}

	return aths, nil
}

func (a *AuthorsQ) Get(ctx context.Context) (modelsdb.AuthorMongo, error) {
	var ath modelsdb.AuthorMongo

	err := a.collection.FindOne(ctx, a.filters).Decode(&ath)
	if err != nil {
		return modelsdb.AuthorMongo{}, err
	}

	return ath, nil
}

func (a *AuthorsQ) FiltersID(id uuid.UUID) *AuthorsQ {
	a.filters["_id"] = id

	return a
}

func (a *AuthorsQ) FiltersName(name string) *AuthorsQ {
	a.filters["name"] = bson.M{
		"$regex":   fmt.Sprintf(".*%s.*", name),
		"$options": "i",
	}

	return a
}

func (a *AuthorsQ) Update(ctx context.Context, fields map[string]any) (modelsdb.AuthorMongo, error) {
	validFields := map[string]bool{
		"name":       true,
		"desc":       true,
		"avatar":     true,
		"email":      true,
		"telegram":   true,
		"twitter":    true,
		"updated_at": true,
	}

	updateFields := bson.M{}
	for key, value := range fields {
		if validFields[key] {
			updateFields[key] = value
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated modelsdb.AuthorMongo

	err := a.collection.FindOneAndUpdate(ctx, a.filters, bson.M{"$set": updateFields}, opts).Decode(&updated)
	if err != nil {
		return modelsdb.AuthorMongo{}, err
	}

	for key, value := range updateFields {
		if _, exists := a.filters[key]; exists {
			a.filters[key] = value
		}
	}

	return updated, nil
}

func (a *AuthorsQ) Limit(limit int64) *AuthorsQ {
	a.limit = limit

	return a
}

func (a *AuthorsQ) Skip(skip int64) *AuthorsQ {
	a.skip = skip

	return a
}

func (a *AuthorsQ) Sort(field string, ascending bool) *AuthorsQ {
	order := 1
	if !ascending {
		order = -1
	}

	a.sort = bson.D{{field, order}}

	return a
}
