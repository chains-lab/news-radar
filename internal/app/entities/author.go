package entities

import (
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
)

type authorsRepo interface {
	Create(author modelsdb.Author) error
	Update(ID uuid.UUID, fields map[string]any) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (modelsdb.Author, error)
}

type Authors struct {
	data authorsRepo
}

func NewAuthors(cfg config.Config) (*Authors, error) {
	data, err := repo.NewAuthors(cfg)
	if err != nil {
		return nil, err
	}

	return &Authors{
		data: data,
	}, nil
}
