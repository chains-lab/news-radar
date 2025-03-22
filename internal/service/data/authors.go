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

type Authors interface {
	Create(ctx context.Context, author models.Author) error
	Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Author, error)
}

type authors struct {
	redis redisdb.Authors
	mongo mongodb.Authors
	neo   neodb.Authors
}

func NewAuthors(cfg config.Config) (Authors, error) {
	mongo, err := mongodb.NewAuthors(cfg.Database.Mongo.URI, cfg.Database.Mongo.Name)
	if err != nil {
		return nil, err
	}
	neo, err := neodb.NewAuthors(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}
	return &authors{
		redis: nil,
		mongo: mongo,
		neo:   neo,
	}, nil
}

func (a *authors) Create(ctx context.Context, author models.Author) error {
	if err := a.neo.Create(ctx, &neodb.Author{
		ID:     author.ID,
		Name:   author.Name,
		Status: author.Status,
	}); err != nil {
		return err
	}

	_, err := a.mongo.New().Insert(ctx, &mongodb.Author{
		ID:        author.ID,
		Name:      author.Name,
		CreatedAt: author.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *authors) Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error {
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

func (a *authors) Delete(ctx context.Context, ID uuid.UUID) error {
	if err := a.neo.Delete(ctx, ID); err != nil {
		return err
	}

	if err := a.mongo.New().FiltersID(ID).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (a *authors) GetByID(ctx context.Context, ID uuid.UUID) (*models.Author, error) {
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

func createModelsAuthor(neo neodb.Author, mongo mongodb.Author) *models.Author {
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
