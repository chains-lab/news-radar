package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
	"github.com/hs-zavet/news-radar/internal/repo/neodb"
)

type usersNeo interface {
	Create(ctx context.Context, user modelsdb.UserNeo) error
	Get(ctx context.Context, id uuid.UUID) (modelsdb.UserNeo, error)
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

func (u *Users) Create(userID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	return u.neo.Create(ctxSync, modelsdb.UserNeo{
		ID: userID,
	})
}

func (u *Users) Get(userID uuid.UUID) (modelsdb.User, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	user, err := u.neo.Get(ctxSync, userID)
	if err != nil {
		return modelsdb.User{}, err
	}

	res := modelsdb.NewUser(user)
	//if err != nil {
	//	return models.User{}, err
	//}

	return res, nil
}
