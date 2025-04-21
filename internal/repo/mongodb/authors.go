package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AuthorsCollection = "authors"
)

type AuthorModel struct {
	ID        uuid.UUID          `json:"_id" bson:"_id"`
	Status    enums.AuthorStatus `json:"status" bson:"status"`
	Name      string             `json:"name" bson:"name"`
	Desc      *string            `json:"desc" bson:"desc"`
	Avatar    *string            `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email     *string            `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string            `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string            `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type AuthorsQ struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection

	filters bson.M
	sort    bson.D
	limit   int64
	skip    int64
}

func NewAuthors(dbName, uri string) (*AuthorsQ, error) {
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

type AuthorInsertInput struct {
	ID        uuid.UUID          `json:"_id" bson:"_id"`
	Status    enums.AuthorStatus `json:"status" bson:"status"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (a *AuthorsQ) Insert(ctx context.Context, input AuthorInsertInput) (AuthorModel, error) {
	author := AuthorModel{
		ID:        input.ID,
		Name:      input.Name,
		CreatedAt: input.CreatedAt,
	}
	author.Status = input.Status

	_, err := a.collection.InsertOne(ctx, author)
	if err != nil {
		return AuthorModel{}, fmt.Errorf("failed to insert author: %w", err)
	}

	return author, nil
}

func (a *AuthorsQ) Delete(ctx context.Context) error {
	_, err := a.collection.DeleteOne(ctx, a.filters)
	if err != nil {
		return fmt.Errorf("failed to delete author: %w", err)
	}

	return nil
}

func (a *AuthorsQ) Count(ctx context.Context) (int64, error) {
	return a.collection.CountDocuments(ctx, a.filters)
}

func (a *AuthorsQ) Select(ctx context.Context) ([]AuthorModel, error) {
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
		return nil, fmt.Errorf("failed to find author: %w", err)
	}

	defer cursor.Close(ctx)

	var aths []AuthorModel

	if err = cursor.All(ctx, &aths); err != nil {
		return nil, fmt.Errorf("failed to decode author: %w", err)
	}

	return aths, nil
}

func (a *AuthorsQ) Get(ctx context.Context) (AuthorModel, error) {
	var ath AuthorModel

	err := a.collection.FindOne(ctx, a.filters).Decode(&ath)
	if err != nil {
		return AuthorModel{}, fmt.Errorf("failed to find author: %w", err)
	}

	return ath, nil
}

func (a *AuthorsQ) FilterID(id uuid.UUID) *AuthorsQ {
	a.filters["_id"] = id

	return a
}

func (a *AuthorsQ) FilterName(name string) *AuthorsQ {
	a.filters["name"] = bson.M{
		"$regex":   fmt.Sprintf(".*%s.*", name),
		"$options": "i",
	}

	return a
}

type AuthorUpdateInput struct {
	Name      *string             `json:"name" bson:"name"`
	Status    *enums.AuthorStatus `json:"status" bson:"status"`
	Desc      *string             `json:"desc" bson:"desc"`
	Avatar    *string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email     *string             `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string             `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string             `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
}

func (a *AuthorsQ) Update(ctx context.Context, input AuthorUpdateInput) (AuthorModel, error) {
	setFields := bson.M{"updated_at": input.UpdatedAt}
	unsetFields := bson.M{}

	// helper: decide to set or unset
	try := func(key string, ptr *string) {
		if ptr == nil {
			return
		}
		if s := strings.TrimSpace(*ptr); s == "" {
			unsetFields[key] = "" // удаляем поле
		} else {
			setFields[key] = s // обновляем
		}
	}

	try("name", input.Name)
	try("desc", input.Desc)
	try("avatar", input.Avatar)
	try("email", input.Email)
	try("telegram", input.Telegram)
	try("twitter", input.Twitter)

	if input.Status != nil {
		setFields["status"] = *input.Status
	}

	if len(setFields) == 1 && len(unsetFields) == 0 {
		return AuthorModel{}, fmt.Errorf("nothing to update")
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}
	if len(unsetFields) > 0 {
		update["$unset"] = unsetFields
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated AuthorModel
	if err := a.collection.FindOneAndUpdate(ctx, a.filters, update, opts).Decode(&updated); err != nil {
		return AuthorModel{}, fmt.Errorf("failed to update author: %w", err)
	}

	// опционально обновляем filters для цепочек
	for k, v := range setFields {
		if _, ok := a.filters[k]; ok {
			a.filters[k] = v
		}
	}
	for k := range unsetFields {
		delete(a.filters, k)
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
