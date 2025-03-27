package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
)

type userRepo interface {
	Create(ctx context.Context, userID uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) (*neodb.UserModels, error)
}

type User struct {
	data userRepo
}

func NewUser(cfg config.Config) (*User, error) {
	repo, err := data.NewUsers(cfg)
	if err != nil {
		return nil, err
	}

	return &User{
		data: repo,
	}, nil
}
