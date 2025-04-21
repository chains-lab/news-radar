package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hs-zavet/news-radar/internal/app/ape"
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

func (a App) CreateTag(ctx context.Context, request CreateTagRequest) (models.Tag, error) {
	CreatedAt := time.Now().UTC()

	tag, err := a.tags.Create(repo.TagCreateInput{
		Name:      request.Name,
		Status:    request.Status,
		Type:      request.Type,
		Color:     request.Color,
		Icon:      request.Icon,
		CreatedAt: CreatedAt,
	})
	if err != nil {
		return models.Tag{}, ape.ErrorTagNameAlreadyTaken
	}

	return models.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Status:    tag.Status,
		Type:      tag.Type,
		Color:     tag.Color,
		Icon:      tag.Icon,
		CreatedAt: tag.CreatedAt,
	}, nil
}

func (a App) DeleteTag(ctx context.Context, id string) error {
	_, err := a.tags.Get(strings.ToLower(id))
	if err != nil {
		return ape.ErrTagNotFound
	}

	err = a.tags.Delete(strings.ToLower(strings.ToLower(id)))
	if err != nil {
		return err
	}

	return nil
}

type UpdateTagRequest struct {
	Name   *string          `json:"name"`
	Status *enums.TagStatus `json:"status"`
	Type   *enums.TagType   `json:"type"`
	Color  *string          `json:"color"`
	Icon   *string          `json:"icon"`
}

func (a App) UpdateTag(ctx context.Context, id string, request UpdateTagRequest) (models.Tag, error) {
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
		if err == nil {
			return models.Tag{}, fmt.Errorf("tag with id %s already exists", *request.Name)
		}
		input.Name = request.Name
	}

	res, err := a.tags.Update(strings.ToLower(id), input)
	if err != nil {
		return models.Tag{}, err
	}

	return models.Tag{
		ID:        res.ID,
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (a App) GetTag(ctx context.Context, id string) (models.Tag, error) {
	res, err := a.tags.Get(strings.ToLower(id))
	if err != nil {
		return models.Tag{}, ape.ErrTagNotFound
	}

	return models.Tag{
		ID:        res.ID,
		Name:      res.Name,
		Status:    res.Status,
		Type:      res.Type,
		Color:     res.Color,
		Icon:      res.Icon,
		CreatedAt: res.CreatedAt,
	}, nil
}
