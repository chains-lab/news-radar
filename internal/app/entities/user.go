package entities

import (
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
)

type userRepo interface {
	Create(userID uuid.UUID) error
	Get(userID uuid.UUID) (modelsdb.User, error)
}

type User struct {
	data userRepo
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
