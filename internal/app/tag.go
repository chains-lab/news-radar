package app

import (
	"context"
	"fmt"
	"time"

	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type CreateTagRequest struct {
	Name   string          `json:"name"`
	Type   enums.TagType   `json:"type"`
	Status enums.TagStatus `json:"status"`
	Color  string          `json:"color"`
	Icon   string          `json:"icon"`
}

func (a App) CreateTag(ctx context.Context, request CreateTagRequest) error {
	CreatedAt := time.Now().UTC()

	return a.tags.Create(repo.TagCreateInput{
		Name:      request.Name,
		Status:    request.Status,
		Type:      request.Type,
		Color:     request.Color,
		Icon:      request.Icon,
		CreatedAt: CreatedAt,
	})
}

func (a App) DeleteTag(ctx context.Context, name string) error {
	return a.tags.Delete(name)
}

type UpdateTagRequest struct {
	Name   *string          `json:"name"`
	Status *enums.TagStatus `json:"status"`
	Type   *enums.TagType   `json:"type"`
	Color  *string          `json:"color"`
	Icon   *string          `json:"icon"`
}

func (a App) UpdateTag(ctx context.Context, name string, request UpdateTagRequest) error {
	input := repo.TagUpdateInput{}

	if request.Status != nil {
		input.Status = request.Status
	}

	if request.Type != nil {
		input.Type = request.Type
	}
	if request.Color != nil {
		input.Color = request.Color
	}
	if request.Icon != nil {
		input.Icon = request.Icon
	}
	if request.Name != nil {
		_, err := a.tags.Get(*request.Name)
		if err != nil {
			return fmt.Errorf("tag with name %s already exists", *request.Name)
		}
		input.Name = request.Name
	}

	return a.tags.Update(name, input)
}

func (a App) GetTag(ctx context.Context, name string) (models.Tag, error) {
	res, err := a.tags.Get(name)
	if err != nil {
		return models.Tag{}, fmt.Errorf("tag with name %s not found", name)
	}

	return models.Tag{
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}
