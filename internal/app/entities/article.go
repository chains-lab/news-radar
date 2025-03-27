package entities

import (
	"context"

	"github.com/google/uuid"
	"github.com/recovery-flow/news-radar/internal/app/models"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/data"
)

type articlesRepo interface {
	Create(ctx context.Context, article models.Article) error
	Update(ctx context.Context, ID uuid.UUID, fields map[string]any) error
	Delete(ctx context.Context, ID uuid.UUID) error

	SetTags(ctx context.Context, ID uuid.UUID, tags []string) error
	AddTag(ctx context.Context, ID uuid.UUID, tag string) error
	DeleteTag(ctx context.Context, ID uuid.UUID, tag string) error

	SetTheme(ctx context.Context, ID uuid.UUID, theme []string) error
	AddTheme(ctx context.Context, ID uuid.UUID, theme string) error
	DeleteTheme(ctx context.Context, ID uuid.UUID, theme string) error

	AddAuthor(ctx context.Context, ID uuid.UUID, author uuid.UUID) error
	DeleteAuthor(ctx context.Context, ID uuid.UUID, author uuid.UUID) error
	SetAuthors(ctx context.Context, ID uuid.UUID, authors []uuid.UUID) error

	GetByID(ctx context.Context, ID uuid.UUID) (*models.Article, error)
}

type Articles struct {
	data articlesRepo
}

func NewArticles(cfg config.Config) (*Articles, error) {
	repo, err := data.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	return &Articles{
		data: repo,
	}, nil
}
