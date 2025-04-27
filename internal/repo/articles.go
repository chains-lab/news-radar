package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo/mongodb"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

const (
	dataCtxTimeAisle = 10 * time.Second
)

type ArticleModel struct {
	ID          uuid.UUID           `json:"id" bson:"_id"`
	Status      enums.ArticleStatus `json:"status" bson:"status"`
	Title       string              `json:"title" bson:"title"`
	Icon        *string             `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc        *string             `json:"desc,omitempty" bson:"desc,omitempty"`
	Content     []content.Section   `json:"content,omitempty" bson:"content,omitempty"`
	PublishedAt *time.Time          `json:"published_at,omitempty" bson:"published_at,omitempty"`
	UpdatedAt   *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
}

type articlesMongoQ interface {
	New() *mongodb.ArticlesQ

	Insert(ctx context.Context, input mongodb.ArticleInsertInput) (mongodb.ArticleModel, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.ArticleModel, error)
	Get(ctx context.Context) (mongodb.ArticleModel, error)

	FilterID(id uuid.UUID) *mongodb.ArticlesQ
	FilterTitle(title string) *mongodb.ArticlesQ
	FilterDate(filters map[string]any, after bool) *mongodb.ArticlesQ
	FilterStatus(status enums.ArticleStatus) *mongodb.ArticlesQ

	Update(ctx context.Context, input mongodb.ArticleUpdateInput) (mongodb.ArticleModel, error)
	DeleteContentSection(
		ctx context.Context,
		index int,
		updatedAt time.Time,
	) error
	UpdateContentSection(
		ctx context.Context,
		section content.Section,
		updatedAt time.Time,
	) error

	Limit(limit int64) *mongodb.ArticlesQ
	Skip(skip int64) *mongodb.ArticlesQ
	Sort(field string, ascending bool) *mongodb.ArticlesQ
}

type articlesNeoQ interface {
	Create(ctx context.Context, input neodb.ArticleInsertInput) (neodb.ArticleModel, error)
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (neodb.ArticleModel, error)

	Update(ctx context.Context, ID uuid.UUID, input neodb.ArticleUpdateInput) (neodb.ArticleModel, error)

	TopicSearch(
		ctx context.Context,
		tag string,
		start, limit int,
	) ([]neodb.ArticleModel, error)

	RecommendByTopic(
		ctx context.Context,
		articleID uuid.UUID,
		limit int,
	) ([]neodb.ArticleModel, error)
}

type ArticlesRepo struct {
	mongo articlesMongoQ
	neo   articlesNeoQ

	hashtag    hashtag
	authorship authorship
}

func NewArticles(cfg config.Config) (*ArticlesRepo, error) {
	mongo, err := mongodb.NewArticles(cfg.Database.Mongo.Name, cfg.Database.Mongo.URI)
	if err != nil {
		return nil, err
	}
	neo, err := neodb.NewArticles(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	hashtagNeo, err := neodb.NewHashtag(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	authorshipNeo, err := neodb.NewAuthorship(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	return &ArticlesRepo{
		mongo: mongo,
		neo:   neo,

		hashtag:    hashtagNeo,
		authorship: authorshipNeo,
	}, nil
}

type ArticleCreateInput struct {
	ID        uuid.UUID           `json:"id" bson:"_id"`
	Title     string              `json:"title" bson:"title"`
	Status    enums.ArticleStatus `json:"status" bson:"status"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

func (a *ArticlesRepo) Create(input ArticleCreateInput) (ArticleModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().Insert(ctxSync, mongodb.ArticleInsertInput{
		ID:        input.ID,
		Status:    input.Status,
		Title:     input.Title,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return ArticleModel{}, err
	}

	_, err = a.neo.Create(ctxSync, neodb.ArticleInsertInput{
		ID:     input.ID,
		Status: input.Status,
	})
	if err != nil {
		return ArticleModel{}, err
	}

	return articleMongoToRepo(mongo), nil
}

type ArticleUpdateInput struct {
	Status      *enums.ArticleStatus `json:"status" bson:"status"`
	PublishedAt *time.Time           `json:"published_at,omitempty" bson:"published_at,omitempty"`
	Title       *string              `json:"title" bson:"title"`
	Icon        *string              `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc        *string              `json:"desc,omitempty" bson:"desc,omitempty"`
}

func (a *ArticlesRepo) Update(ID uuid.UUID, input ArticleUpdateInput) (ArticleModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	updatedAt := time.Now().UTC()

	var neoInput neodb.ArticleUpdateInput
	var mongoInput mongodb.ArticleUpdateInput

	if status := input.Status; status != nil {
		neoInput.Status = status
		mongoInput.Status = status
	}

	if input.PublishedAt != nil {
		neoInput.PublishedAt = input.PublishedAt
		mongoInput.PublishedAt = input.PublishedAt
	}

	if input.Title != nil {
		mongoInput.Title = input.Title
	}
	if input.Icon != nil {
		mongoInput.Icon = input.Icon
	}
	if input.Desc != nil {
		mongoInput.Desc = input.Desc
	}

	mongoInput.UpdatedAt = updatedAt

	_, err := a.neo.Update(ctxSync, ID, neoInput)
	if err != nil {
		return ArticleModel{}, err
	}

	mongo, err := a.mongo.New().FilterID(ID).Update(ctxSync, mongoInput)
	if err != nil {
		return ArticleModel{}, err
	}

	return articleMongoToRepo(mongo), nil
}

func (a *ArticlesRepo) DeleteContentSection(ID uuid.UUID, sectionID int) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	updatedAt := time.Now().UTC()

	err := a.mongo.New().FilterID(ID).DeleteContentSection(ctxSync, sectionID, updatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) Delete(ID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.neo.Delete(ctxSync, ID)
	if err != nil {
		return err
	}

	err = a.mongo.FilterID(ID).Delete(ctxSync)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) GetByID(ID uuid.UUID) (ArticleModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongoRes, err := a.mongo.FilterID(ID).Get(ctxSync)
	if err != nil {
		return ArticleModel{}, err
	}

	res := articleMongoToRepo(mongoRes)

	return res, nil
}

func (a *ArticlesRepo) RecommendByTopic(
	ctx context.Context,
	articleID uuid.UUID,
	limit int,
) ([]ArticleModel, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	res, err := a.neo.RecommendByTopic(ctxSync, articleID, limit)
	if err != nil {
		return nil, err
	}

	var articles []ArticleModel
	for _, item := range res {
		mongoRes, err := a.mongo.FilterID(item.ID).Get(ctxSync)
		if err != nil {
			return nil, err
		}

		model := articleMongoToRepo(mongoRes)

		articles = append(articles, model)
	}

	return articles, nil
}

func (a *ArticlesRepo) TopicSearch(
	ctx context.Context,
	tag string,
	start, limit int,
) ([]ArticleModel, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	res, err := a.neo.TopicSearch(ctxSync, tag, start, limit)
	if err != nil {
		return nil, err
	}

	var articles []ArticleModel
	for _, item := range res {
		mongoRes, err := a.mongo.FilterID(item.ID).Get(ctxSync)
		if err != nil {
			return nil, err
		}

		model := articleMongoToRepo(mongoRes)

		articles = append(articles, model)
	}

	return articles, nil
}

func articleMongoToRepo(mongo mongodb.ArticleModel) ArticleModel {
	res := ArticleModel{
		ID:        mongo.ID,
		Status:    mongo.Status,
		Title:     mongo.Title,
		CreatedAt: mongo.CreatedAt,
	}

	if mongo.Icon != nil {
		res.Icon = mongo.Icon
	}

	if mongo.Desc != nil {
		res.Desc = mongo.Desc
	}

	if mongo.PublishedAt != nil {
		res.PublishedAt = mongo.PublishedAt
	}

	if mongo.UpdatedAt != nil {
		res.UpdatedAt = mongo.UpdatedAt
	}

	if mongo.Content != nil {
		res.Content = mongo.Content
	}

	return res
}
