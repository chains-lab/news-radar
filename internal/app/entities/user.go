package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data/neodb"
	"github.com/recovery-flow/news-radar/internal/repo"
)

type UserRepo interface {
	Create(ctx context.Context, userID uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) (*neodb.UserModels, error)

	AddLike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveLike(ctx context.Context, userID, articleID uuid.UUID) error

	AddDislike(ctx context.Context, userID, articleID uuid.UUID) error
	RemoveDislike(ctx context.Context, userID, articleID uuid.UUID) error

	AddRepost(ctx context.Context, userID, articleID uuid.UUID) error
}

type User struct {
	data UserRepo
}

func NewUser(cfg config.Config) (*User, error) {
	repo, err := repo.NewUsers(cfg)
	if err != nil {
		return nil, err
	}

	return &User{
		data: repo,
	}, nil
}
