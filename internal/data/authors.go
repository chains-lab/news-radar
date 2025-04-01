package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/data/models"
	"github.com/hs-zavet/news-radar/internal/data/mongodb"
	"github.com/hs-zavet/news-radar/internal/data/neodb"
)

type authorsRedis interface {
}

type authorsMongo interface {
	New() *mongodb.AuthorsQ

	Insert(ctx context.Context, author mongodb.AuthorModel) (mongodb.AuthorModel, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.AuthorModel, error)
	Get(ctx context.Context) (mongodb.AuthorModel, error)

	FiltersID(id uuid.UUID) *mongodb.AuthorsQ
	FiltersName(name string) *mongodb.AuthorsQ

	Update(ctx context.Context, fields map[string]any) (mongodb.AuthorModel, error)

	Limit(limit int64) *mongodb.AuthorsQ
	Skip(skip int64) *mongodb.AuthorsQ
	Sort(field string, ascending bool) *mongodb.AuthorsQ
}

type authorsNeo interface {
	Create(ctx context.Context, author neodb.AuthorModel) error
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (neodb.AuthorModel, error)

	UpdateName(ctx context.Context, ID uuid.UUID, name string) error
	UpdateStatus(ctx context.Context, ID uuid.UUID, status string) error
}

type Authors struct {
	redis authorsRedis
	mongo authorsMongo
	neo   authorsNeo
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

func (a *Authors) Create(author models.Author) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if err := a.neo.Create(ctxSync, neodb.AuthorModel{
		ID:     author.ID,
		Name:   author.Name,
		Status: author.Status,
	}); err != nil {
		return err
	}

	_, err := a.mongo.New().Insert(ctxSync, mongodb.AuthorModel{
		ID:        author.ID,
		Name:      author.Name,
		CreatedAt: author.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *Authors) Update(ID uuid.UUID, fields map[string]any) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if status, ok := fields["status"].(string); ok {
		if err := a.neo.UpdateStatus(ctxSync, ID, status); err != nil {
			return err
		}
	}

	if name, ok := fields["name"].(string); ok {
		if err := a.neo.UpdateName(ctxSync, ID, name); err != nil {
			return err
		}
	}

	_, err := a.mongo.New().FiltersID(ID).Update(ctxSync, fields)
	if err != nil {
		return err
	}

	return nil
}

func (a *Authors) Delete(ID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if err := a.neo.Delete(ctxSync, ID); err != nil {
		return err
	}

	if err := a.mongo.New().FiltersID(ID).Delete(ctxSync); err != nil {
		return err
	}

	return nil
}

func (a *Authors) GetByID(ID uuid.UUID) (models.Author, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().FiltersID(ID).Get(ctxSync)
	if err != nil {
		return models.Author{}, err
	}

	neo, err := a.neo.GetByID(ctxSync, ID)
	if err != nil {
		return models.Author{}, err
	}

	res, err := models.AuthorsCreateModel(mongo, neo)
	if err != nil {
		return models.Author{}, err
	}

	return res, nil
}
