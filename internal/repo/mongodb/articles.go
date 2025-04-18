package mongodb

import (
	"context"
	"errors"
	"fmt"
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
	ID        uuid.UUID `json:"__id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (a *ArticlesQ) Insert(ctx context.Context, input ArticleInsertInput) error {
	_, err := a.collection.InsertOne(ctx, ArticleModel{
		ID:        input.ID,
		Status:    enums.ArticleStatusActive,
		Title:     input.Title,
		Icon:      nil,
		Desc:      nil,
		Content:   nil,
		UpdatedAt: nil,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
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
		return nil, fmt.Errorf("failed to select ArticlesQ: %w", err)
	}

	defer cursor.Close(ctx)

	var arts []ArticleModel
	if err = cursor.All(ctx, &arts); err != nil {
		return nil, fmt.Errorf("failed to decode ArticlesQ: %w", err)
	}

	return arts, nil
}

func (a *ArticlesQ) Get(ctx context.Context) (ArticleModel, error) {
	var art ArticleModel

	err := a.collection.FindOne(ctx, a.filters).Decode(&art)
	if err != nil {
		return ArticleModel{}, err
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

func (a *ArticlesQ) applyUpdate(ctx context.Context, updatedAt time.Time, update bson.M) (ArticleModel, error) {
	// гарантируем, что updated_at будет записан
	setUpdated := func(m bson.M) {
		if _, ok := m["updated_at"]; !ok {
			m["updated_at"] = updatedAt
		}
	}

	// если в апдейте есть $set/$currentDate, добавляем туда updated_at
	if set, ok := update["$set"].(bson.M); ok {
		setUpdated(set)
	} else if cd, ok := update["$currentDate"].(bson.M); ok {
		setUpdated(cd)
	} else {
		// вообще нет $set или $currentDate → создаём $set
		update["$set"] = bson.M{"updated_at": updatedAt}
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated ArticleModel
	if err := a.collection.FindOneAndUpdate(ctx, a.filters, update, opts).Decode(&updated); err != nil {
		return ArticleModel{}, err
	}

	// синхронизируем filters, чтобы можно было звать методы цепочкой
	for k, v := range update {
		if k != "$set" {
			continue
		}
		if setMap, ok := v.(bson.M); ok {
			for field, val := range setMap {
				a.filters[field] = val
			}
		}
	}

	return updated, nil
}

// applySet — сахар над applyUpdate для простых $set обновлений.
func (a *ArticlesQ) applySet(ctx context.Context, updatedAt time.Time, fields bson.M) (ArticleModel, error) {
	return a.applyUpdate(ctx, updatedAt, bson.M{"$set": fields})
}

func (a *ArticlesQ) UpdateStatus(ctx context.Context, updatedAt time.Time, status enums.ArticleStatus) (ArticleModel, error) {
	return a.applySet(ctx, updatedAt, bson.M{"status": status})
}

func (a *ArticlesQ) UpdateTitle(ctx context.Context, updatedAt time.Time, title string) (ArticleModel, error) {
	return a.applySet(ctx, updatedAt, bson.M{"title": title})
}

// UpdateIcon позволяет как заменить иконку, так и сбросить её, передав nil.
func (a *ArticlesQ) UpdateIcon(ctx context.Context, updatedAt time.Time, icon *string) (ArticleModel, error) {
	return a.applySet(ctx, updatedAt, bson.M{"icon": icon})
}

func (a *ArticlesQ) UpdateDesc(ctx context.Context, updatedAt time.Time, desc *string) (ArticleModel, error) {
	return a.applySet(ctx, updatedAt, bson.M{"desc": desc})
}

// UpdateContent изменяет, удаляет или добавляет секцию контента по правилам:
//
//	— если section «пустая» (нет media, audio и текста) и index < len(content) → удалить;
//	— если не пустая и index < len(content) → заменить существующую;
//	— если не пустая и index >= len(content) → push в конец;
//	— если пустая и index >= len(content) → только updated_at.
//
// Возвращает итоговую версию статьи.
func (a *ArticlesQ) UpdateContent(ctx context.Context, updatedAt time.Time, index int, section content.Section) (ArticleModel, error) {
	// 1) Тянем действующую статью, чтобы знать длину content
	var art ArticleModel
	if err := a.collection.FindOne(ctx, a.filters).Decode(&art); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ArticleModel{}, fmt.Errorf("article not found")
		}
		return ArticleModel{}, fmt.Errorf("failed to load article: %w", err)
	}

	// 2) Проверяем «пустая» ли секция
	isEmpty := func(s content.Section) bool {
		return s.Media == nil && len(s.Text) == 0 && s.Audio == nil
	}

	// 3) Формируем update
	var update bson.M

	switch {
	case isEmpty(section) && index < len(art.Content):
		update = bson.M{
			"$unset":       bson.M{fmt.Sprintf("content.%d", index): 1},
			"$pull":        bson.M{"content": nil},
			"$currentDate": bson.M{"updated_at": updatedAt},
		}
	case !isEmpty(section) && index < len(art.Content):
		update = bson.M{
			"$set":         bson.M{fmt.Sprintf("content.%d", index): section},
			"$currentDate": bson.M{"updated_at": updatedAt},
		}
	case !isEmpty(section) && index >= len(art.Content):
		update = bson.M{
			"$push":        bson.M{"content": section},
			"$currentDate": bson.M{"updated_at": updatedAt},
		}
	default:
		update = bson.M{"$set": bson.M{"updated_at": updatedAt}}
	}

	// 4) Применяем и получаем обновлённую статью
	return a.applyUpdate(ctx, updatedAt, update)
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
