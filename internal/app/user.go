package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type CreateUserRequest struct {
	ID uuid.UUID `json:"id"`
}

func (a App) CreateUser(ctx context.Context, request CreateUserRequest) error {
	return a.users.Create(repo.UserCreateInput{
		ID: request.ID,
	})
}
