package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

type UserModel struct {
	ID uuid.UUID `json:"id"`
}

type usersNeo interface {
	Create(ctx context.Context, input neodb.UserCreateInput) error
	Get(ctx context.Context, id uuid.UUID) (neodb.UserModel, error)
}

type Users struct {
	neo usersNeo
}

func NewUsers(cfg config.Config) (*Users, error) {
	neo, err := neodb.NewUsers(cfg.Database.Neo4j.URI, cfg.Database.Neo4j.Username, cfg.Database.Neo4j.Password)
	if err != nil {
		return nil, err
	}

	return &Users{
		neo: neo,
	}, nil
}

type UserCreateInput struct {
	ID uuid.UUID `json:"id"`
}

func (u *Users) Create(input UserCreateInput) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.neo.Create(ctxSync, neodb.UserCreateInput{
		ID: input.ID,
	})
}

func (u *Users) Get(userID uuid.UUID) (UserModel, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	user, err := u.neo.Get(ctxSync, userID)
	if err != nil {
		return UserModel{}, err
	}

	res := NewUser(user)
	//if err != nil {
	//	return models.User{}, err
	//}

	return res, nil
}

func NewUser(neo neodb.UserModel) UserModel {
	return UserModel{
		ID: neo.ID,
	}
}
