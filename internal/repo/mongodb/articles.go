package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ArticlesCollection = "articles"
)

type ArticlesQ struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection

	filters bson.M
	sort    bson.D
	limit   int64
	skip    int64
}

func NewArticles(cfg config.Config) (*ArticlesQ, error) {
	clientOptions := options.Client().ApplyURI(cfg.Database.Mongo.URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	database := client.Database(cfg.Database.Mongo.Name)
	coll := database.Collection(ArticlesCollection)

	return &ArticlesQ{
		client:     client,
		database:   database,
		collection: coll,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}, nil
}

func (a *ArticlesQ) New() *ArticlesQ {
	return &ArticlesQ{
		client:     a.client,
		database:   a.database,
		collection: a.collection,
		filters:    bson.M{},
		sort:       bson.D{},
		limit:      0,
		skip:       0,
	}
}

func (a *ArticlesQ) Insert(ctx context.Context, article modelsdb.ArticleMongo) (modelsdb.ArticleMongo, error) {
	_, err := a.collection.InsertOne(ctx, article)
	if err != nil {
		return modelsdb.ArticleMongo{}, err
	}

	return article, nil
}

func (a *ArticlesQ) Delete(ctx context.Context) error {
	_, err := a.collection.DeleteOne(ctx, a.filters)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesQ) Count(ctx context.Context) (int64, error) {
	return a.collection.CountDocuments(ctx, a.filters)
}

func (a *ArticlesQ) Select(ctx context.Context) ([]modelsdb.ArticleMongo, error) {
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
		return nil, fmt.Errorf("failed to select ArticlesQ: %w", err)
	}

	defer cursor.Close(ctx)

	var arts []modelsdb.ArticleMongo
	if err = cursor.All(ctx, &arts); err != nil {
		return nil, fmt.Errorf("failed to decode ArticlesQ: %w", err)
	}

	return arts, nil
}

func (a *ArticlesQ) Get(ctx context.Context) (modelsdb.ArticleMongo, error) {
	var art modelsdb.ArticleMongo

	err := a.collection.FindOne(ctx, a.filters).Decode(&art)
	if err != nil {
		return modelsdb.ArticleMongo{}, err
	}

	return art, nil
}

func (a *ArticlesQ) FilterID(id uuid.UUID) *ArticlesQ {
	a.filters["_id"] = id

	return a
}

func (a *ArticlesQ) FilterTitle(title string) *ArticlesQ {
	a.filters["title"] = bson.M{
		"$regex":   fmt.Sprintf(".*%s.*", title),
		"$options": "i",
	}

	return a
}

func (a *ArticlesQ) FilterDate(filters map[string]any, after bool) *ArticlesQ {
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

func (a *ArticlesQ) Update(ctx context.Context, fields map[string]any) (modelsdb.ArticleMongo, error) {
	validFields := map[string]bool{
		"title":      true,
		"icon":       true,
		"desc":       true,
		"content":    true,
		"likes":      true,
		"reposts":    true,
		"updated_at": true,
	}
	updateFields := bson.M{}
	for key, value := range fields {
		if validFields[key] {
			updateFields[key] = value
		}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated modelsdb.ArticleMongo
	err := a.collection.FindOneAndUpdate(ctx, a.filters, bson.M{"$set": updateFields}, opts).Decode(&updated)
	if err != nil {
		return modelsdb.ArticleMongo{}, fmt.Errorf("failed to update article: %w", err)
	}

	return updated, nil
}

func (a *ArticlesQ) Limit(limit int64) *ArticlesQ {
	a.limit = limit

	return a
}

func (a *ArticlesQ) Skip(skip int64) *ArticlesQ {
	a.skip = skip

	return a
}

func (a *ArticlesQ) Sort(field string, ascending bool) *ArticlesQ {
	order := 1
	if !ascending {
		order = -1
	}

	a.sort = bson.D{{Key: field, Value: order}}

	return a
}
