package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type userRepo interface {
	Create(userID repo.UserCreateInput) error
	Get(userID uuid.UUID) (repo.UserModel, error)
}

type User struct {
	data userRepo
}

func NewUser(cfg config.Config) (*User, error) {
	data, err := repo.NewUsers(cfg)
	if err != nil {
		return nil, err
	}

	return &User{
		data: data,
	}, nil
}

type CreateUserRequest struct {
	ID uuid.UUID `json:"id"`
}

func (u *User) CreateUser(ctx context.Context, request CreateUserRequest) error {
	return u.data.Create(repo.UserCreateInput{
		ID: request.ID,
	})
}
