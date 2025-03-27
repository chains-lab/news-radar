package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type usersNeo interface {
	Create(ctx context.Context, user neodb.UserModels) error
	Get(ctx context.Context, id uuid.UUID) (*neodb.UserModels, error)
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

func (u *Users) Create(ctx context.Context, userID uuid.UUID) error {
	return u.neo.Create(ctx, neodb.UserModels{
		ID: userID,
	})
}

func (u *Users) Get(ctx context.Context, userID uuid.UUID) (*neodb.UserModels, error) {
	return u.neo.Get(ctx, userID)
}
