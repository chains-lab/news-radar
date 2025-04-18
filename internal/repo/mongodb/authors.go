package mongodb

import (
	"context"
	"fmt"
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

func (a *AuthorsQ) Insert(ctx context.Context, input AuthorInsertInput) error {
	stmt := AuthorModel{
		ID:        input.ID,
		Name:      input.Name,
		CreatedAt: input.CreatedAt,
	}
	stmt.Status = input.Status

	_, err := a.collection.InsertOne(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
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
		return nil, err
	}

	defer cursor.Close(ctx)

	var aths []AuthorModel

	if err = cursor.All(ctx, &aths); err != nil {
		return nil, err
	}

	return aths, nil
}

func (a *AuthorsQ) Get(ctx context.Context) (AuthorModel, error) {
	var ath AuthorModel

	err := a.collection.FindOne(ctx, a.filters).Decode(&ath)
	if err != nil {
		return AuthorModel{}, err
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

func (a *AuthorsQ) applyUpdates(ctx context.Context, fields bson.M, updatedAt time.Time) (AuthorModel, error) {
	fields["updated_at"] = updatedAt

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated AuthorModel

	err := a.collection.
		FindOneAndUpdate(ctx, a.filters, bson.M{"$set": fields}, opts).
		Decode(&updated)
	if err != nil {
		return AuthorModel{}, err
	}

	for k, v := range fields {
		a.filters[k] = v
	}

	return updated, nil
}

func (a *AuthorsQ) UpdateName(ctx context.Context, name *string, updatedAt time.Time) (AuthorModel, error) {
	if name == nil {
		return AuthorModel{}, fmt.Errorf("name is required")
	}
	return a.applyUpdates(ctx, bson.M{"name": *name}, updatedAt)
}

func (a *AuthorsQ) UpdateDescription(ctx context.Context, desc *string, updatedAt time.Time) (AuthorModel, error) {
	if desc == nil {
		return AuthorModel{}, fmt.Errorf("desc is required")
	}
	return a.applyUpdates(ctx, bson.M{"desc": *desc}, updatedAt)
}

func (a *AuthorsQ) UpdateAvatar(ctx context.Context, avatar *string, updatedAt time.Time) (AuthorModel, error) {
	return a.applyUpdates(ctx, bson.M{"avatar": avatar}, updatedAt)
}

func (a *AuthorsQ) UpdateStatus(ctx context.Context, status enums.AuthorStatus, updatedAt time.Time) (AuthorModel, error) {
	if status == "" {
		return AuthorModel{}, fmt.Errorf("status is required")
	}
	return a.applyUpdates(ctx, bson.M{"status": status}, updatedAt)
}

func (a *AuthorsQ) UpdateContactInfo(
	ctx context.Context,
	email, telegram, twitter *string,
	updatedAt time.Time,
) (AuthorModel, error) {
	fields := bson.M{}
	if email != nil {
		fields["email"] = *email
	}
	if telegram != nil {
		fields["telegram"] = *telegram
	}
	if twitter != nil {
		fields["twitter"] = *twitter
	}
	if len(fields) == 0 {
		return AuthorModel{}, fmt.Errorf("nothing to update in contacts")
	}
	return a.applyUpdates(ctx, fields, updatedAt)
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
