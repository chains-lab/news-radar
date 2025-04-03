package app

import (
	"context"
	"fmt"
	"time"

	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type CreateTagRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func (a App) CreateTag(ctx context.Context, request CreateTagRequest) error {
	CreatedAt := time.Now().UTC()

	tType, err := models.ParseTagType(request.Type)
	if err != nil {
		return err
	}

	return a.tags.Create(repo.TagCreateInput{
		Name:      request.Name,
		Status:    string(models.TagStatusInactive),
		Type:      string(tType),
		Color:     request.Color,
		Icon:      request.Icon,
		CreatedAt: CreatedAt,
	})
}

func (a App) DeleteTag(ctx context.Context, name string) error {
	return a.tags.Delete(name)
}

type UpdateTagRequest struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
	Type   *string `json:"type"`
	Color  *string `json:"color"`
	Icon   *string `json:"icon"`
}

func (a App) UpdateTag(ctx context.Context, name string, request UpdateTagRequest) error {
	input := repo.TagUpdateInput{}

	if request.Status != nil {
		_, err := models.ParseTagStatus(*request.Status)
		if err != nil {
			return err
		}
		input.Status = request.Status
	}

	if request.Type != nil {
		_, err := models.ParseTagType(*request.Type)
		if err != nil {
			return err
		}
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

func (a App) Get(ctx context.Context, name string) (repo.TagModel, error) {
	return a.tags.Get(name)
}
