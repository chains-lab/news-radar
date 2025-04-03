package repo

import (
	"context"

	"github.com/google/uuid"
)

type authorship interface {
	Create(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error
	Delete(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error

	SetForArticle(ctx context.Context, ID uuid.UUID, author []uuid.UUID) error

	GetForArticle(ctx context.Context, ID uuid.UUID) ([]uuid.UUID, error)
	GetForAuthor(ctx context.Context, ID uuid.UUID) ([]uuid.UUID, error)
}

func (a *ArticlesRepo) SetAuthors(articleID uuid.UUID, authors []uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.authorship.SetForArticle(ctxSync, articleID, authors)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) AddAuthor(articleID uuid.UUID, authorID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.authorship.Create(ctxSync, articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) DeleteAuthor(articleID uuid.UUID, authorID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.authorship.Delete(ctxSync, articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) GetAuthors(articleID uuid.UUID) ([]uuid.UUID, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	authors, err := a.authorship.GetForArticle(ctxSync, articleID)
	if err != nil {
		return nil, err
	}

	return authors, nil
}
