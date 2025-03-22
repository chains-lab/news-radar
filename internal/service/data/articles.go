package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/data/mongodb"
	"github.com/recovery-flow/news-radar/internal/service/data/neodb"
	"github.com/recovery-flow/news-radar/internal/service/data/redisdb"
	"github.com/recovery-flow/news-radar/internal/service/models"
)

type Article interface {
	Create(ctx context.Context, article models.Article) error
	Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error
	Delete(ctx context.Context, ID uuid.UUID) error

	SetTags(ctx context.Context, ID uuid.UUID, tags []string) error
	AddTag(ctx context.Context, ID uuid.UUID, tag string) error
	DeleteTag(ctx context.Context, ID uuid.UUID, tag string) error

	SetTheme(ctx context.Context, ID uuid.UUID, theme []string) error
	AddTheme(ctx context.Context, ID uuid.UUID, theme string) error
	DeleteTheme(ctx context.Context, ID uuid.UUID, theme string) error

	AddAuthor(ctx context.Context, ID uuid.UUID, author uuid.UUID) error
	DeleteAuthor(ctx context.Context, ID uuid.UUID, author uuid.UUID) error
	SetAuthors(ctx context.Context, ID uuid.UUID, authors []uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error)
}

type articles struct {
	redis redisdb.Articles
	mongo mongodb.Articles
	neo   neodb.Articles

	hashtag    neodb.Hashtag
	authorship neodb.Authorship
	about      neodb.About
}

func NewArticles(cfg config.Config) (Article, error) {
	mongo, err := mongodb.NewArticles(cfg.Database.Mongo.URI, cfg.Database.Mongo.Name)
	if err != nil {
		return nil, err
	}
	neo, err := neodb.NewArticles(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	hashtag, err := neodb.NewHashtag(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	authorship, err := neodb.NewAuthorship(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	about, err := neodb.NewAbout(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	return &articles{
		redis: nil,
		mongo: mongo,
		neo:   neo,

		hashtag:    hashtag,
		authorship: authorship,
		about:      about,
	}, nil
}

func (a *articles) Create(ctx context.Context, article models.Article) error {
	err := a.neo.Create(ctx, &neodb.Article{
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
		Status:    article.Status,
	})
	if err != nil {
		return err
	}

	_, err = a.mongo.New().Insert(ctx, &mongodb.Article{
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

	return nil
}

func (a *articles) Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error {
	if _, ok := fields["status"]; ok {
		st, err := models.ParseArticleStatus(fields["status"].(string))
		if err != nil {
			return err
		}
		err = a.neo.UpdateStatus(ctx, ID, st)
		if err != nil {
			return err
		}
	}

	_, err := a.mongo.FilterID(ID).Update(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) Delete(ctx context.Context, ID uuid.UUID) error {
	err := a.neo.Delete(ctx, ID)
	if err != nil {
		return err
	}

	err = a.mongo.FilterID(ID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error) {
	mongoRes, err := a.mongo.FilterID(ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	neoRes, err := a.neo.Get(ctx, ID)
	if err != nil {
		return nil, err
	}

	res := createModelArticle(*neoRes, *mongoRes)

	return &res, nil
}

func (a *articles) SetTags(ctx context.Context, articleID uuid.UUID, tags []string) error {
	err := a.hashtag.SetForArticle(ctx, articleID, tags)
	if err != nil {
		return err
	}
	return nil
}

func (a *articles) AddTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	err := a.hashtag.Create(ctx, articleID, tag)
	if err != nil {
		return err
	}
	return nil
}

func (a *articles) DeleteTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	err := a.hashtag.Delete(ctx, articleID, tag)
	if err != nil {
		return err
	}
	return nil
}

func (a *articles) SetTheme(ctx context.Context, articleID uuid.UUID, theme []string) error {
	err := a.about.SetForArticle(ctx, articleID, theme)
	if err != nil {
		return err
	}
	return nil
}

func (a *articles) AddTheme(ctx context.Context, articleID uuid.UUID, theme string) error {
	err := a.about.Create(ctx, articleID, theme)
	if err != nil {
		return err
	}
	return nil
}

func (a *articles) DeleteTheme(ctx context.Context, articleID uuid.UUID, theme string) error {
	err := a.about.Delete(ctx, articleID, theme)
	if err != nil {
		return err
	}
	return nil
}

func (a *articles) SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error {
	err := a.authorship.SetForArticle(ctx, articleID, authors)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) AddAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	err := a.authorship.Create(ctx, articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a *articles) DeleteAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	err := a.authorship.Delete(ctx, articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func createModelArticle(neo neodb.Article, mongo mongodb.Article) models.Article {
	return models.Article{
		ID:        mongo.ID,
		Title:     mongo.Title,
		Icon:      mongo.Icon,
		Desc:      mongo.Desc,
		Content:   mongo.Content,
		Likes:     mongo.Likes,
		Reposts:   mongo.Reposts,
		Status:    neo.Status,
		CreatedAt: mongo.CreatedAt,
	}
}
