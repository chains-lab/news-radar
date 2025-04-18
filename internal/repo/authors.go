package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo/mongodb"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

type AuthorModel struct {
	ID        uuid.UUID          `json:"id" bson:"id"`
	Title     string             `json:"title" bson:"title"`
	Status    enums.AuthorStatus `json:"status" bson:"status"`
	Desc      *string            `json:"desc" bson:"desc"`
	Avatar    *string            `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email     *string            `json:"email,omitempty" bson:"email,omitempty"`
	Telegram  *string            `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter   *string            `json:"twitter,omitempty" bson:"twitter,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type authorsMongo interface {
	New() *mongodb.AuthorsQ

	Insert(ctx context.Context, input mongodb.AuthorInsertInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.AuthorModel, error)
	Get(ctx context.Context) (mongodb.AuthorModel, error)

	FilterID(id uuid.UUID) *mongodb.AuthorsQ
	FilterName(name string) *mongodb.AuthorsQ

	UpdateStatus(ctx context.Context, status enums.AuthorStatus, updatedAt time.Time) (mongodb.AuthorModel, error)
	UpdateDescription(ctx context.Context, desc *string, updatedAt time.Time) (mongodb.AuthorModel, error)
	UpdateAvatar(ctx context.Context, avatar *string, updatedAt time.Time) (mongodb.AuthorModel, error)
	UpdateContactInfo(
		ctx context.Context,
		email, telegram, twitter *string,
		updatedAt time.Time,
	) (mongodb.AuthorModel, error)

	Limit(limit int64) *mongodb.AuthorsQ
	Skip(skip int64) *mongodb.AuthorsQ
	Sort(field string, ascending bool) *mongodb.AuthorsQ
}

type authorsNeo interface {
	Create(ctx context.Context, input neodb.AuthorCreateInput) error
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (neodb.AuthorModel, error)

	UpdateStatus(ctx context.Context, ID uuid.UUID, status enums.AuthorStatus) error
}

type Authors struct {
	redis interface{}
	mongo authorsMongo
	neo   authorsNeo
}

func NewAuthors(cfg config.Config) (*Authors, error) {
	mongo, err := mongodb.NewAuthors(cfg.Database.Mongo.Name, cfg.Database.Mongo.URI)
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

type AuthorCreateInput struct {
	ID        uuid.UUID          `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Status    enums.AuthorStatus `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (a *Authors) Create(input AuthorCreateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.mongo.New().Insert(ctxSync, mongodb.AuthorInsertInput{
		ID:        input.ID,
		Name:      input.Name,
		Status:    input.Status,
		CreatedAt: input.CreatedAt,
	})

	if err = a.neo.Create(ctxSync, neodb.AuthorCreateInput{
		ID:     input.ID,
		Status: input.Status,
	}); err != nil {
		return err
	}

	return nil
}

func (a *Authors) UpdateStatus(ID uuid.UUID, status enums.AuthorStatus, updatedAt time.Time) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if err := a.neo.UpdateStatus(ctxSync, ID, status); err != nil {
		return AuthorModel{}, err
	}

	mongo, err := a.mongo.New().FilterID(ID).UpdateStatus(ctxSync, status, updatedAt)
	if err != nil {
		return AuthorModel{}, err
	}

	res, err := authorsCreateModel(mongo, neodb.AuthorModel{
		ID:     ID,
		Status: status,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	return res, nil
}

func (a *Authors) UpdateDescription(ID uuid.UUID, desc *string, updatedAt time.Time) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().FilterID(ID).UpdateDescription(ctxSync, desc, updatedAt)
	if err != nil {
		return AuthorModel{}, err
	}

	res, err := authorsCreateModel(mongo, neodb.AuthorModel{
		ID:     ID,
		Status: mongo.Status,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	return res, nil
}

func (a *Authors) UpdateAvatar(ID uuid.UUID, avatar *string, updatedAt time.Time) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().FilterID(ID).UpdateAvatar(ctxSync, avatar, updatedAt)
	if err != nil {
		return AuthorModel{}, err
	}

	res, err := authorsCreateModel(mongo, neodb.AuthorModel{
		ID:     ID,
		Status: mongo.Status,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	return res, nil
}

func (a *Authors) UpdateContactInfo(
	ID uuid.UUID, email, telegram, twitter *string, updatedAt time.Time) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().FilterID(ID).UpdateContactInfo(ctxSync, email, telegram, twitter, updatedAt)
	if err != nil {
		return AuthorModel{}, err
	}

	res, err := authorsCreateModel(mongo, neodb.AuthorModel{
		ID:     ID,
		Status: mongo.Status,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	return res, nil
}

func (a *Authors) Delete(ID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	if err := a.neo.Delete(ctxSync, ID); err != nil {
		return err
	}

	if err := a.mongo.New().FilterID(ID).Delete(ctxSync); err != nil {
		return err
	}

	return nil
}

func (a *Authors) GetByID(ID uuid.UUID) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return AuthorModel{}, err
	}

	neo, err := a.neo.GetByID(ctxSync, ID)
	if err != nil {
		return AuthorModel{}, err
	}

	res, err := authorsCreateModel(mongo, neo)
	if err != nil {
		return AuthorModel{}, err
	}

	return res, nil
}

func authorsCreateModel(mongo mongodb.AuthorModel, neo neodb.AuthorModel) (AuthorModel, error) {
	if mongo.ID != neo.ID {
		return AuthorModel{}, fmt.Errorf("mongo and neo IDs do not match")
	}

	res := AuthorModel{
		ID:        mongo.ID,
		Title:     mongo.Name,
		Status:    neo.Status,
		CreatedAt: mongo.CreatedAt,
	}

	if mongo.Desc != nil {
		res.Desc = mongo.Desc
	}
	if mongo.Avatar != nil {
		res.Avatar = mongo.Avatar
	}
	if mongo.Email != nil {
		res.Email = mongo.Email
	}
	if mongo.Telegram != nil {
		res.Telegram = mongo.Telegram
	}
	if mongo.Twitter != nil {
		res.Twitter = mongo.Twitter
	}
	if mongo.UpdatedAt != nil {
		res.UpdatedAt = mongo.UpdatedAt
	}

	return res, nil
}
