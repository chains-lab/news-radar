package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/internal/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ArticlesCollection = "articles"
)

type ArticleModel struct {
	ID        uuid.UUID           `json:"_id" bson:"_id"`
	Status    enums.ArticleStatus `json:"status" bson:"status"`
	Title     string              `json:"title" bson:"title"`
	Icon      *string             `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc      *string             `json:"desc,omitempty" bson:"desc,omitempty"`
	Content   []content.Section   `json:"content,omitempty" bson:"content,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

type ArticlesQ struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection

	filters bson.M
	sort    bson.D
	limit   int64
	skip    int64
}

func NewArticles(dbname string, uri string) (*ArticlesQ, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	database := client.Database(dbname)
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

type ArticleInsertInput struct {
	ID        uuid.UUID           `json:"__id" bson:"_id"`
	Title     string              `json:"title" bson:"title"`
	Status    enums.ArticleStatus `json:"status" bson:"status"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

func (a *ArticlesQ) Insert(ctx context.Context, input ArticleInsertInput) (ArticleModel, error) {
	article := ArticleModel{
		ID:        input.ID,
		Status:    input.Status,
		Title:     input.Title,
		Icon:      nil,
		Desc:      nil,
		Content:   nil,
		UpdatedAt: nil,
		CreatedAt: input.CreatedAt,
	}

	_, err := a.collection.InsertOne(ctx, article)
	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to insert article: %w", err)
	}

	return article, nil
}

func (a *ArticlesQ) Delete(ctx context.Context) error {
	_, err := a.collection.DeleteOne(ctx, a.filters)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	return nil
}

func (a *ArticlesQ) Count(ctx context.Context) (int64, error) {
	return a.collection.CountDocuments(ctx, a.filters)
}

func (a *ArticlesQ) Select(ctx context.Context) ([]ArticleModel, error) {
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
		return nil, fmt.Errorf("failed to select article: %w", err)
	}

	defer cursor.Close(ctx)

	var arts []ArticleModel
	if err = cursor.All(ctx, &arts); err != nil {
		return nil, fmt.Errorf("failed to convert article to models: %w", err)
	}

	return arts, nil
}

func (a *ArticlesQ) Get(ctx context.Context) (ArticleModel, error) {
	var art ArticleModel

	err := a.collection.FindOne(ctx, a.filters).Decode(&art)
	if err != nil {
		return ArticleModel{}, fmt.Errorf("failed to get article: %w", err)
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

func (a *ArticlesQ) FilterStatus(status enums.ArticleStatus) *ArticlesQ {
	a.filters["status"] = string(status)

	return a
}

type ArticleUpdateInput struct {
	Status    *enums.ArticleStatus `json:"status" bson:"status"`
	Title     *string              `json:"title,omitempty" bson:"title,omitempty"`
	Icon      *string              `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc      *string              `json:"desc,omitempty" bson:"desc,omitempty"`
	UpdatedAt time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (a *ArticlesQ) Update(ctx context.Context, input ArticleUpdateInput) (ArticleModel, error) {
	setFields := bson.M{"updated_at": input.UpdatedAt}
	unsetFields := bson.M{}

	try := func(key string, ptr *string) {
		if ptr == nil {
			return
		}
		if s := strings.TrimSpace(*ptr); s == "" {
			unsetFields[key] = ""
		} else {
			setFields[key] = s
		}
	}

	try("title", input.Title)
	try("icon", input.Icon)
	try("desc", input.Desc)

	if input.Status != nil {
		setFields["status"] = *input.Status
	}

	if len(setFields) == 1 && len(unsetFields) == 0 {
		return ArticleModel{}, fmt.Errorf("nothing to update")
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}
	if len(unsetFields) > 0 {
		update["$unset"] = unsetFields
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated ArticleModel
	if err := a.collection.FindOneAndUpdate(ctx, a.filters, update, opts).Decode(&updated); err != nil {
		return ArticleModel{}, fmt.Errorf("failed to update article: %w", err)
	}

	return updated, nil
}

func (a *ArticlesQ) UpdateContent(
	ctx context.Context,
	index int,
	section content.Section,
	updatedAt time.Time,
) (ArticleModel, error) {
	// 1. Считали статью
	var article ArticleModel
	if err := a.collection.FindOne(ctx, a.filters).Decode(&article); err != nil {
		return ArticleModel{}, fmt.Errorf("failed to find article: %w", err)
	}

	// 2. Определяем, пустая ли секция
	isEmpty := func(s content.Section) bool {
		return s.Media == nil &&
			len(s.Text) == 0 &&
			s.Audio == nil
	}

	// 3. Собираем update
	var update bson.M

	switch {
	// Удаление: пересобираем слайс без элемента с индексом
	case isEmpty(section) && index < len(article.Content):
		newContent := make([]content.Section, 0, len(article.Content)-1)
		for i, sec := range article.Content {
			if i != index {
				newContent = append(newContent, sec)
			}
		}
		update = bson.M{
			"$set": bson.M{
				"content":    newContent,
				"updated_at": updatedAt,
			},
		}

	// Перезапись существующей
	case !isEmpty(section) && index < len(article.Content):
		update = bson.M{
			"$set": bson.M{
				fmt.Sprintf("content.%d", index): section,
				"updated_at":                     updatedAt,
			},
		}

	// Добавление в конец
	case !isEmpty(section) && index >= len(article.Content):
		update = bson.M{
			"$push": bson.M{"content": section},
			"$set":  bson.M{"updated_at": updatedAt},
		}

	// Пустая секция вне диапазона — только updated_at
	default:
		update = bson.M{
			"$set": bson.M{"updated_at": updatedAt},
		}
	}

	// 4. Применяем с опцией возврата обновлённого документа
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated ArticleModel
	if err := a.collection.
		FindOneAndUpdate(ctx, a.filters, update, opts).
		Decode(&updated); err != nil {
		return ArticleModel{}, fmt.Errorf(
			"failed to update content at index %d: %w", index, err,
		)
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
