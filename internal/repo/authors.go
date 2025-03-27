package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/mongodb"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type AuthorsRedis interface {
}

type AuthorsMongo interface {
	New() *mongodb.AuthorsQ

	Insert(ctx context.Context, author *mongodb.AuthorModel) (*mongodb.AuthorModel, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.AuthorModel, error)
	Get(ctx context.Context) (*mongodb.AuthorModel, error)

	FiltersID(id uuid.UUID) *mongodb.AuthorsQ
	FiltersName(name string) *mongodb.AuthorsQ

	Update(ctx context.Context, fields map[string]any) (*mongodb.AuthorModel, error)

	Limit(limit int64) *mongodb.AuthorsQ
	Skip(skip int64) *mongodb.AuthorsQ
	Sort(field string, ascending bool) *mongodb.AuthorsQ
}

type AuthorsNeo interface {
	Create(ctx context.Context, author *neodb.AuthorModel) error
	Delete(ctx context.Context, ID uuid.UUID) error

	UpdateName(ctx context.Context, ID uuid.UUID, name string) error
	UpdateStatus(ctx context.Context, ID uuid.UUID, status models.AuthorStatus) error

	GetByID(ctx context.Context, ID uuid.UUID) (*neodb.AuthorModel, error)
}

type Authors struct {
	redis AuthorsRedis
	mongo AuthorsMongo
	neo   AuthorsNeo
}

func NewAuthors(cfg config.Config) (*Authors, error) {
	mongo, err := mongodb.NewAuthors(cfg.Database.Mongo.URI, cfg.Database.Mongo.Name)
	if err != nil {
		return nil, err
	}
	neo, err := neodb.NewAuthors(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	return &Authors{
		redis: nil,
		mongo: mongo,
		neo:   neo,
	}, nil
}

func (a *Authors) Create(ctx context.Context, author models.Author) error {
	if err := a.neo.Create(ctx, &neodb.AuthorModel{
		ID:     author.ID,
		Name:   author.Name,
		Status: author.Status,
	}); err != nil {
		return err
	}

	_, err := a.mongo.New().Insert(ctx, &mongodb.AuthorModel{
		ID:        author.ID,
		Name:      author.Name,
		CreatedAt: author.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *Authors) Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error {
	if _, ok := fields["status"]; ok {
		status, err := models.ParseAuthorStatus(fields["status"].(string))
		if err != nil {
			return err
		}
		if err := a.neo.UpdateStatus(ctx, ID, status); err != nil {
			return err
		}
	}

	if _, ok := fields["name"]; ok {
		if err := a.neo.UpdateName(ctx, ID, fields["name"].(string)); err != nil {
			return err
		}
	}

	_, err := a.mongo.New().FiltersID(ID).Update(ctx, fields)
	if err != nil {
		return err
	}

	return nil
}

func (a *Authors) Delete(ctx context.Context, ID uuid.UUID) error {
	if err := a.neo.Delete(ctx, ID); err != nil {
		return err
	}

	if err := a.mongo.New().FiltersID(ID).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Authors) GetByID(ctx context.Context, ID uuid.UUID) (*models.Author, error) {
	mongo, err := a.mongo.New().FiltersID(ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	neo, err := a.neo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	res := createModelsAuthor(*neo, *mongo)

	return res, nil
}

func createModelsAuthor(neo neodb.AuthorModel, mongo mongodb.AuthorModel) *models.Author {
	return &models.Author{
		ID:        mongo.ID,
		Name:      mongo.Name,
		CreatedAt: mongo.CreatedAt,
		Desc:      mongo.Desc,
		Avatar:    mongo.Avatar,
		Email:     mongo.Email,
		Telegram:  mongo.Telegram,
		Twitter:   mongo.Twitter,
		UpdatedAt: mongo.UpdatedAt,
		Status:    neo.Status,
	}
}
