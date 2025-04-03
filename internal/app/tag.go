package domain

import (
	"context"
	"time"

	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type tagsRepo interface {
	Create(input repo.TagCreateInput) error
	Delete(name string) error
	Update(name string, input repo.TagUpdateInput) error
	Get(name string) (repo.TagModel, error)
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

type CreateTagRequest struct {
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (t *Tags) CreateTag(ctx context.Context, name string, request CreateTagRequest) error {
	CreatedAt := time.Now().UTC()

	return t.data.Create(repo.TagCreateInput{
		Name:      name,
		Status:    string(models.TagStatusInactive),
		Color:     request.Color,
		Icon:      request.Icon,
		CreatedAt: CreatedAt,
	})
}

func (t *Tags) DeleteTag(ctx context.Context, name string) error {
	return t.data.Delete(name)
}

type UpdateTagRequest struct {
	Color  *string `json:"color"`
	Icon   *string `json:"icon"`
	Status *string `json:"status"`
}

func (t *Tags) UpdateTag(ctx context.Context, name string, request UpdateTagRequest) error {
	return t.data.Update(name, repo.TagUpdateInput{
		Color:  request.Color,
		Icon:   request.Icon,
		Status: request.Status,
	})
}

func (t *Tags) Get(ctx context.Context, name string) (repo.TagModel, error) {
	return t.data.Get(name)
}
