package repo

import (
	"context"
	"fmt"
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
	ID        uuid.UUID           `json:"id" bson:"_id"`
	Status    enums.ArticleStatus `json:"status" bson:"status"`
	Title     string              `json:"title" bson:"title"`
	Icon      *string             `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc      *string             `json:"desc,omitempty" bson:"desc,omitempty"`
	Content   []content.Section   `json:"content,omitempty" bson:"content,omitempty"`
	UpdatedAt *time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

type articlesMongoQ interface {
	New() *mongodb.ArticlesQ

	Insert(ctx context.Context, input mongodb.ArticleInsertInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.ArticleModel, error)
	Get(ctx context.Context) (mongodb.ArticleModel, error)

	FilterID(id uuid.UUID) *mongodb.ArticlesQ
	FilterTitle(title string) *mongodb.ArticlesQ
	FilterDate(filters map[string]any, after bool) *mongodb.ArticlesQ

	Update(ctx context.Context, input mongodb.ArticleUpdateInput) (mongodb.ArticleModel, error)

	Limit(limit int64) *mongodb.ArticlesQ
	Skip(skip int64) *mongodb.ArticlesQ
	Sort(field string, ascending bool) *mongodb.ArticlesQ
}

type articlesNeoQ interface {
	Create(ctx context.Context, input neodb.ArticleInsertInput) error
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (neodb.ArticleModel, error)

	UpdateStatus(ctx context.Context, ID uuid.UUID, status enums.ArticleStatus) error
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

func (a *ArticlesRepo) Create(input ArticleCreateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.mongo.New().Insert(ctxSync, mongodb.ArticleInsertInput{
		ID:        input.ID,
		Title:     input.Title,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return err
	}

	err = a.neo.Create(ctxSync, neodb.ArticleInsertInput{
		ID:     input.ID,
		Status: input.Status,
	})
	if err != nil {
		return err
	}

	return nil
}

type ArticleUpdateInput struct {
	Status    *enums.ArticleStatus `json:"status,omitempty" bson:"status,omitempty"`
	Title     *string              `json:"title,omitempty" bson:"title,omitempty"`
	Icon      *string              `json:"icon,omitempty" bson:"icon,omitempty"`
	Desc      *string              `json:"desc,omitempty" bson:"desc,omitempty"`
	Content   []content.Section    `json:"content,omitempty" bson:"content,omitempty"`
	UpdatedAt time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (a *ArticlesRepo) Update(ID uuid.UUID, input ArticleUpdateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if input.Status != nil {
		err := a.neo.UpdateStatus(ctxSync, ID, *input.Status)
		if err != nil {
			return err
		}
	}

	_, err := a.mongo.FilterID(ID).Update(ctxSync, mongodb.ArticleUpdateInput{
		Status:    input.Status,
		Title:     input.Title,
		Icon:      input.Icon,
		Desc:      input.Desc,
		Content:   input.Content,
		UpdatedAt: input.UpdatedAt,
	})
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

	neoRes, err := a.neo.GetByID(ctxSync, ID)
	if err != nil {
		return ArticleModel{}, err
	}

	res, err := CreateArticleModel(mongoRes, neoRes)
	if err != nil {
		return ArticleModel{}, err
	}

	return res, nil
}

func CreateArticleModel(mongo mongodb.ArticleModel, neo neodb.ArticleModel) (ArticleModel, error) {
	if mongo.ID != neo.ID {
		return ArticleModel{}, fmt.Errorf("mongo and neo IDs do not match")
	}

	return ArticleModel{
		ID:        mongo.ID,
		Title:     mongo.Title,
		Icon:      mongo.Icon,
		Desc:      mongo.Desc,
		Content:   mongo.Content,
		Status:    neo.Status,
		UpdatedAt: mongo.UpdatedAt,
		CreatedAt: mongo.CreatedAt,
	}, nil
}
