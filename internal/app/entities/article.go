package entities

import (
	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
)

type articlesRepo interface {
	Create(article modelsdb.Article) error
	Update(ID uuid.UUID, fields map[string]any) error
	Delete(ID uuid.UUID) error

	SetTags(ID uuid.UUID, tags []string) error
	AddTag(ID uuid.UUID, tag string) error
	DeleteTag(ID uuid.UUID, tag string) error

	AddAuthor(ID uuid.UUID, author uuid.UUID) error
	DeleteAuthor(ID uuid.UUID, author uuid.UUID) error
	SetAuthors(ID uuid.UUID, authors []uuid.UUID) error

	GetByID(ID uuid.UUID) (modelsdb.Article, error)
}

type Articles struct {
	data articlesRepo
}

func NewArticles(cfg config.Config) (*Articles, error) {
	data, err := repo.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	return &Articles{
		data: data,
	}, nil
}
