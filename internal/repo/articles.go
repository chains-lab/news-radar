package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
	"github.com/hs-zavet/news-radar/internal/repo/mongodb"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

const (
	dataCtxTimeAisle = 10 * time.Second
)

type articlesMongo interface {
	New() *mongodb.ArticlesQ

	Insert(ctx context.Context, article modelsdb.ArticleMongo) (modelsdb.ArticleMongo, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]modelsdb.ArticleMongo, error)
	Get(ctx context.Context) (modelsdb.ArticleMongo, error)

	FilterID(id uuid.UUID) *mongodb.ArticlesQ
	FilterTitle(title string) *mongodb.ArticlesQ
	FilterDate(filters map[string]any, after bool) *mongodb.ArticlesQ

	Update(ctx context.Context, fields map[string]any) (modelsdb.ArticleMongo, error)

	Limit(limit int64) *mongodb.ArticlesQ
	Skip(skip int64) *mongodb.ArticlesQ
	Sort(field string, ascending bool) *mongodb.ArticlesQ
}

type articlesNeo interface {
	Create(ctx context.Context, article modelsdb.ArticleNeo) error
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (modelsdb.ArticleNeo, error)

	UpdateStatus(ctx context.Context, ID uuid.UUID, status string) error
}

type ArticlesRepo struct {
	mongo articlesMongo
	neo   articlesNeo

	hashtag    hashtag
	authorship authorship
}

func NewArticles(cfg config.Config) (*ArticlesRepo, error) {
	mongo, err := mongodb.NewArticles(cfg)
	if err != nil {
		return nil, err
	}
	neo, err := neodb.NewArticles(cfg)
	if err != nil {
		return nil, err
	}
	hashtagNeo, err := neodb.NewHashtag(cfg)
	if err != nil {
		return nil, err
	}
	authorshipNeo, err := neodb.NewAuthorship(cfg)
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

func (a *ArticlesRepo) Create(article modelsdb.Article) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	_, err := a.mongo.New().Insert(ctxSync, modelsdb.ArticleMongo{
		ID:        article.ID,
		Title:     article.Title,
		Icon:      article.Icon,
		Desc:      article.Desc,
		Content:   article.Content,
		Likes:     article.Likes,
		Reposts:   article.Reposts,
		CreatedAt: article.CreatedAt,
	})
	if err != nil {
		return err
	}

	err = a.neo.Create(ctxSync, modelsdb.ArticleNeo{
		ID:     article.ID,
		Status: article.Status,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) Update(ID uuid.UUID, fields map[string]any) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if _, ok := fields["status"]; ok {
		st := fields["status"].(string)

		err := a.neo.UpdateStatus(ctxSync, ID, st)
		if err != nil {
			return err
		}
	}

	_, err := a.mongo.FilterID(ID).Update(ctxSync, fields)
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

func (a *ArticlesRepo) GetByID(ID uuid.UUID) (modelsdb.Article, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongoRes, err := a.mongo.FilterID(ID).Get(ctxSync)
	if err != nil {
		return modelsdb.Article{}, err
	}

	neoRes, err := a.neo.GetByID(ctxSync, ID)
	if err != nil {
		return modelsdb.Article{}, err
	}

	res, err := modelsdb.CreateArticleModel(mongoRes, neoRes)
	if err != nil {
		return modelsdb.Article{}, err
	}

	return res, nil
}
