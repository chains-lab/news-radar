package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type authorsRepo interface {
	Create(author repo.AuthorCreateInput) error
	Update(ID uuid.UUID, input repo.AuthorUpdateInput) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (repo.AuthorModel, error)
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

type CreateAuthorRequest struct {
	Name     *string `json:"name" bson:"name"`
	Desc     *string `json:"desc" bson:"desc"`
	Avatar   *string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email    *string `json:"email,omitempty" bson:"email,omitempty"`
	Telegram *string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter  *string `json:"twitter,omitempty" bson:"twitter,omitempty"`
}

func (a *Authors) CreateAuthor(ctx context.Context, request CreateAuthorRequest) error {
	AuthorID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.data.Create(repo.AuthorCreateInput{
		ID:        AuthorID,
		Name:      *request.Name,
		Status:    string(models.AuthorStatusInactive),
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
	Name     *string `json:"name" bson:"name"`
	Status   *string `json:"status" bson:"status"`
	Desc     *string `json:"desc" bson:"desc"`
	Avatar   *string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email    *string `json:"email,omitempty" bson:"email,omitempty"`
	Telegram *string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Twitter  *string `json:"twitter,omitempty" bson:"twitter,omitempty"`
}

func (a *Authors) UpdateAuthor(ctx context.Context, authorID uuid.UUID, request UpdateAuthorRequest) error {
	UpdatedAt := time.Now().UTC()

	_, err := models.ParseAuthorStatus(*request.Status)
	if err != nil {
		return err
	}
	return a.data.Update(authorID, repo.AuthorUpdateInput{
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

func (a *Authors) DeleteAuthor(ctx context.Context, authorID uuid.UUID) error {
	return a.data.Delete(authorID)
}

func (a *Authors) GetByID(ctx context.Context, authorID uuid.UUID) (repo.AuthorModel, error) {
	return a.data.GetByID(authorID)
}
