package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type CreateAuthorRequest struct {
	Name     string  `json:"name" bson:"name"`
	Desc     *string `json:"desc" bson:"desc"`
	Avatar   *string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email    *string `json:"email,omitempty" bson:"email,omitempty"`
	Telegram *string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter  *string `json:"twitter,omitempty" bson:"twitter,omitempty"`
}

func (a App) CreateAuthor(ctx context.Context, request CreateAuthorRequest) error {
	AuthorID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.authors.Create(repo.AuthorCreateInput{
		ID:        AuthorID,
		Name:      request.Name,
		Status:    enums.AuthorStatusActive,
		Desc:      request.Desc,
		Avatar:    request.Avatar,
		Email:     request.Email,
		Telegram:  request.Telegram,
		Twitter:   request.Twitter,
		CreatedAt: CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

type UpdateAuthorRequest struct {
	Name     *string             `json:"name" bson:"name"`
	Status   *enums.AuthorStatus `json:"status" bson:"status"`
	Desc     *string             `json:"desc" bson:"desc"`
	Avatar   *string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email    *string             `json:"email,omitempty" bson:"email,omitempty"`
	Telegram *string             `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter  *string             `json:"twitter,omitempty" bson:"twitter,omitempty"`
}

func (a App) UpdateAuthor(ctx context.Context, authorID uuid.UUID, request UpdateAuthorRequest) error {
	UpdatedAt := time.Now().UTC()

	return a.authors.Update(authorID, repo.AuthorUpdateInput{
		Name:      request.Name,
		Status:    request.Status,
		Desc:      request.Desc,
		Avatar:    request.Avatar,
		Email:     request.Email,
		Telegram:  request.Telegram,
		Twitter:   request.Twitter,
		UpdatedAt: UpdatedAt,
	})
}

func (a App) DeleteAuthor(ctx context.Context, authorID uuid.UUID) error {
	return a.authors.Delete(authorID)
}

func (a App) GetAuthorByID(ctx context.Context, authorID uuid.UUID) (models.Author, error) {
	res, err := a.authors.GetByID(authorID)
	if err != nil {
		return models.Author{}, err
	}

	return models.Author{
		ID:        res.ID,
		Name:      res.Name,
		Status:    res.Status,
		Desc:      res.Desc,
		Avatar:    res.Avatar,
		Email:     res.Email,
		Telegram:  res.Telegram,
		Twitter:   res.Twitter,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
