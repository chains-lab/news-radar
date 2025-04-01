package data

import (
	"context"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/data/neodb"
)

type hashtag interface {
	Create(ctx context.Context, articleID uuid.UUID, tag string) error
	Delete(ctx context.Context, articleID uuid.UUID, tag string) error

	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]neodb.TagModels, error)

	SetForArticle(ctx context.Context, articleID uuid.UUID, tags []string) error
}

func (a *ArticlesRepo) SetTags(articleID uuid.UUID, tags []string) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.hashtag.SetForArticle(ctxSync, articleID, tags)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) AddTag(articleID uuid.UUID, tag string) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.hashtag.Create(ctxSync, articleID, tag)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArticlesRepo) DeleteTag(articleID uuid.UUID, tag string) error {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	err := a.hashtag.Delete(ctxSync, articleID, tag)
	if err != nil {
		return err
	}

	return nil
}
