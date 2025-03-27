package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/mongodb"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type ArticlesMongo interface {
	New() *mongodb.ArticlesQ

	Insert(ctx context.Context, article *mongodb.ArticleModel) (*mongodb.ArticleModel, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.ArticleModel, error)
	Get(ctx context.Context) (*mongodb.ArticleModel, error)

	FilterID(id uuid.UUID) *mongodb.ArticlesQ
	FilterTitle(title string) *mongodb.ArticlesQ
	FilterDate(filters map[string]any, after bool) *mongodb.ArticlesQ

	Update(ctx context.Context, fields map[string]any) (*mongodb.ArticleModel, error)

	Limit(limit int64) *mongodb.ArticlesQ
	Skip(skip int64) *mongodb.ArticlesQ
	Sort(field string, ascending bool) *mongodb.ArticlesQ
}

type ArticlesNeo interface {
	Create(ctx context.Context, article *neodb.ArticleModel) error
	Delete(ctx context.Context, ID uuid.UUID) error
	Get(ctx context.Context, ID uuid.UUID) (*neodb.ArticleModel, error)

	UpdateStatus(ctx context.Context, ID uuid.UUID, status models.ArticleStatus) error
}

type Hashtag interface {
	Create(ctx context.Context, articleID uuid.UUID, tag string) error
	Delete(ctx context.Context, articleID uuid.UUID, tag string) error

	SetForArticle(ctx context.Context, articleID uuid.UUID, tags []string) error
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]*models.Tag, error)
}

type Authorship interface {
	Create(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error
	Delete(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error

	SetForArticle(ctx context.Context, ID uuid.UUID, author []uuid.UUID) error
	GetForArticle(ctx context.Context, ID uuid.UUID) ([]uuid.UUID, error)

	GetForAuthor(ctx context.Context, ID uuid.UUID) ([]uuid.UUID, error)
}

type About interface {
	Create(ctx context.Context, articleID uuid.UUID, theme string) error
	Delete(ctx context.Context, articleID uuid.UUID, theme string) error

	SetForArticle(ctx context.Context, articleID uuid.UUID, themes []string) error
	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]*neodb.ThemeModels, error)
}

type ArticlesRepo struct {
	mongo ArticlesMongo
	neo   ArticlesNeo

	hashtag    Hashtag
	authorship Authorship
	about      About
}

func NewArticles(cfg config.Config) (*ArticlesRepo, error) {
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
	return &ArticlesRepo{
		mongo: mongo,
		neo:   neo,

		hashtag:    hashtag,
		authorship: authorship,
		about:      about,
	}, nil
}

func (a *ArticlesRepo) Create(ctx context.Context, article models.Article) error {
	err := a.neo.Create(ctx, &neodb.ArticleModel{
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
		Status:    article.Status,
	})
	if err != nil {
		return err
	}

	_, err = a.mongo.New().Insert(ctx, &mongodb.ArticleModel{
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

func (a *ArticlesRepo) Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error {
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

func (a *ArticlesRepo) Delete(ctx context.Context, ID uuid.UUID) error {
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

func (a *ArticlesRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error) {
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

func (a *ArticlesRepo) SetTags(ctx context.Context, articleID uuid.UUID, tags []string) error {
	err := a.hashtag.SetForArticle(ctx, articleID, tags)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepo) AddTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	err := a.hashtag.Create(ctx, articleID, tag)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepo) DeleteTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	err := a.hashtag.Delete(ctx, articleID, tag)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepo) SetTheme(ctx context.Context, articleID uuid.UUID, theme []string) error {
	err := a.about.SetForArticle(ctx, articleID, theme)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepo) AddTheme(ctx context.Context, articleID uuid.UUID, theme string) error {
	err := a.about.Create(ctx, articleID, theme)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepo) DeleteTheme(ctx context.Context, articleID uuid.UUID, theme string) error {
	err := a.about.Delete(ctx, articleID, theme)
	if err != nil {
		return err
	}
	return nil
}

func (a *ArticlesRepo) SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error {
	err := a.authorship.SetForArticle(ctx, articleID, authors)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) AddAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	err := a.authorship.Create(ctx, articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) DeleteAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	err := a.authorship.Delete(ctx, articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func createModelArticle(neo neodb.ArticleModel, mongo mongodb.ArticleModel) models.Article {
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
