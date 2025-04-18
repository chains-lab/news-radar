package repo

import (
	"context"

	"github.com/google/uuid"
)

type hashtag interface {
	Create(ctx context.Context, articleID uuid.UUID, tag string) error
	Delete(ctx context.Context, articleID uuid.UUID, tag string) error

	GetForArticle(ctx context.Context, articleID uuid.UUID) ([]string, error)
	GetArticlesForTag(ctx context.Context, tag string) ([]uuid.UUID, error)

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

func (a *ArticlesRepo) GetTags(articleID uuid.UUID) ([]string, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	tags, err := a.hashtag.GetForArticle(ctxSync, articleID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (a *ArticlesRepo) GetArticlesForTag(name string) ([]uuid.UUID, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	articles, err := a.hashtag.GetArticlesForTag(ctxSync, name)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
