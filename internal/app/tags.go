package app

import (
	"context"
	"strings"
	"time"

	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/chains-lab/news-radar/internal/app/models"
	"github.com/chains-lab/news-radar/internal/enums"
	"github.com/chains-lab/news-radar/internal/repo"
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

	res := models.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Status:    tag.Status,
		Type:      tag.Type,
		Color:     tag.Color,
		Icon:      tag.Icon,
		CreatedAt: tag.CreatedAt,
	}
	if tag.UpdatedAt != nil {
		res.UpdatedAt = tag.UpdatedAt
	}

	return res, nil
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

	updated := false

	if request.Status != nil {
		input.Status = request.Status
		updated = true
	}
	if request.Type != nil {
		input.Type = request.Type
		updated = true
	}
	if request.Color != nil {
		input.Color = request.Color
		updated = true
	}
	if request.Icon != nil {
		input.Icon = request.Icon
		updated = true
	}
	if request.Name != nil {
		_, err := a.tags.Get(*request.Name)
		if err == nil {
			return models.Tag{}, ape.ErrorTagNameAlreadyTaken
		}
		input.Name = request.Name
		updated = true
	}

	if !updated {
		return a.GetTag(ctx, id)
	}

	_, err := a.tags.Get(strings.ToLower(id))
	if err != nil {
		return models.Tag{}, ape.ErrTagNotFound
	}

	tag, err := a.tags.Update(strings.ToLower(id), input)
	if err != nil {
		return models.Tag{}, err
	}

	res := models.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Status:    tag.Status,
		Type:      tag.Type,
		Color:     tag.Color,
		Icon:      tag.Icon,
		CreatedAt: tag.CreatedAt,
	}
	if tag.UpdatedAt != nil {
		res.UpdatedAt = tag.UpdatedAt
	}

	return res, nil
}

func (a App) GetTag(ctx context.Context, id string) (models.Tag, error) {
	tag, err := a.tags.Get(strings.ToLower(id))
	if err != nil {
		return models.Tag{}, ape.ErrTagNotFound
	}

	res := models.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		Status:    tag.Status,
		Type:      tag.Type,
		Color:     tag.Color,
		Icon:      tag.Icon,
		CreatedAt: tag.CreatedAt,
	}
	if tag.UpdatedAt != nil {
		res.UpdatedAt = tag.UpdatedAt
	}

	return res, nil
}
