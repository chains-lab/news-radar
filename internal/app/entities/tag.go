package entities

import (
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
	"github.com/hs-zavet/news-radar/internal/repo/modelsdb"
)

type tagsRepo interface {
	Create(tag modelsdb.Tag) error
	Delete(name string) error
	Update(name string, fields map[string]any) error
	Get(name string) (modelsdb.Tag, error)
}

type Tags struct {
	data tagsRepo
}

func NewTags(cfg config.Config) (*Tags, error) {
	data, err := repo.NewTags(cfg)
	if err != nil {
		return nil, err
	}

	return &Tags{
		data: data,
	}, nil
}
