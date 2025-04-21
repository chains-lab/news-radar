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
	Name      string             `json:"name" bson:"name"`
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

	Insert(ctx context.Context, input mongodb.AuthorInsertInput) (mongodb.AuthorModel, error)
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int64, error)
	Select(ctx context.Context) ([]mongodb.AuthorModel, error)
	Get(ctx context.Context) (mongodb.AuthorModel, error)

	FilterID(id uuid.UUID) *mongodb.AuthorsQ
	FilterName(name string) *mongodb.AuthorsQ

	Update(ctx context.Context, input mongodb.AuthorUpdateInput) (mongodb.AuthorModel, error)

	Limit(limit int64) *mongodb.AuthorsQ
	Skip(skip int64) *mongodb.AuthorsQ
	Sort(field string, ascending bool) *mongodb.AuthorsQ
}

type authorsNeo interface {
	Create(ctx context.Context, input neodb.AuthorCreateInput) (neodb.AuthorModel, error)
	Delete(ctx context.Context, ID uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (neodb.AuthorModel, error)

	Update(ctx context.Context, id uuid.UUID, input neodb.AuthorUpdateInput) (neodb.AuthorModel, error)
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

func (a *Authors) Create(input AuthorCreateInput) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	mongo, err := a.mongo.New().Insert(ctxSync, mongodb.AuthorInsertInput{
		ID:        input.ID,
		Name:      input.Name,
		Status:    input.Status,
		CreatedAt: input.CreatedAt,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	_, err = a.neo.Create(ctxSync, neodb.AuthorCreateInput{
		ID:     input.ID,
		Status: input.Status,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	return AuthorModel{
		ID:        mongo.ID,
		Name:      mongo.Name,
		Status:    mongo.Status,
		CreatedAt: mongo.CreatedAt,
	}, nil
}

type AuthorUpdateInput struct {
	Name     *string
	Status   *enums.AuthorStatus
	Desc     *string
	Avatar   *string
	Email    *string
	Telegram *string
	Twitter  *string
}

func (a *Authors) Update(ID uuid.UUID, input AuthorUpdateInput) (AuthorModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	updatedAt := time.Now().UTC()

	var mongoInput mongodb.AuthorUpdateInput

	if input.Name != nil {
		mongoInput.Name = input.Name
	}
	if input.Desc != nil {
		mongoInput.Desc = input.Desc
	}
	if input.Avatar != nil {
		mongoInput.Avatar = input.Avatar
	}
	if input.Email != nil {
		mongoInput.Email = input.Email
	}
	if input.Telegram != nil {
		mongoInput.Telegram = input.Telegram
	}
	if input.Twitter != nil {
		mongoInput.Twitter = input.Twitter
	}
	if input.Status != nil {
		_, err := a.neo.Update(ctxSync, ID, neodb.AuthorUpdateInput{
			Status: input.Status,
		})
		if err != nil {
			return AuthorModel{}, err
		}
		mongoInput.Status = input.Status
	}

	mongo, err := a.mongo.New().FilterID(ID).Update(ctxSync, mongodb.AuthorUpdateInput{
		Name:      input.Name,
		Desc:      input.Desc,
		Avatar:    input.Avatar,
		Email:     input.Email,
		Telegram:  input.Telegram,
		Twitter:   input.Twitter,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return AuthorModel{}, err
	}

	return AuthorModel{
		ID:        mongo.ID,
		Name:      mongo.Name,
		Status:    mongo.Status,
		Desc:      mongo.Desc,
		Avatar:    mongo.Avatar,
		Email:     mongo.Email,
		Telegram:  mongo.Telegram,
		Twitter:   mongo.Twitter,
		UpdatedAt: mongo.UpdatedAt,
		CreatedAt: mongo.CreatedAt,
	}, nil
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
		Name:      mongo.Name,
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
