package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/data/models"
	"github.com/hs-zavet/news-radar/internal/data/neodb"
)

type usersNeo interface {
	Create(ctx context.Context, user neodb.UserModels) error
	Get(ctx context.Context, id uuid.UUID) (neodb.UserModels, error)
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

	return u.neo.Create(ctxSync, neodb.UserModels{
		ID: userID,
	})
}

func (u *Users) Get(userID uuid.UUID) (models.User, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	user, err := u.neo.Get(ctxSync, userID)
	if err != nil {
		return models.User{}, err
	}

	res := models.NewUser(user)
	//if err != nil {
	//	return models.User{}, err
	//}

	return res, nil
}
