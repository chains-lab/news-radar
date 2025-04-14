package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type CreateArticleRequest struct {
	Title string `json:"title"`
}

func (a App) CreateArticle(ctx context.Context, request CreateArticleRequest) (models.Article, error) {
	ArticleID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.articles.Create(repo.ArticleCreateInput{
		ID:        ArticleID,
		Title:     request.Title,
		Status:    enums.ArticleStatusInactive,
		CreatedAt: CreatedAt,
	})
	if err != nil {
		return models.Article{}, err
	}

	res, err := a.articles.GetByID(ArticleID)
	if err != nil {
		return models.Article{}, err
	}
	return models.Article{
		ID:        res.ID,
		Title:     res.Title,
		Icon:      nil,
		Desc:      nil,
		Content:   nil,
		Status:    enums.ArticleStatusInactive,
		UpdatedAt: nil,
		CreatedAt: res.CreatedAt,

		Authors: nil,
		Tags:    nil,
	}, nil
}

type UpdateArticleRequest struct {
	Title   *string              `json:"title,omitempty"`
	Status  *enums.ArticleStatus `json:"status,omitempty"`
	Icon    *string              `json:"icon,omitempty"`
	Desc    *string              `json:"desc,omitempty"`
	Content []content.Section    `json:"content,omitempty"`
	Likes   *int                 `json:"likes,omitempty"`
	Reposts *int                 `json:"reposts,omitempty"`
}

func (a App) UpdateArticle(ctx context.Context, articleID uuid.UUID, request UpdateArticleRequest) (models.Article, error) {
	UpdatedAt := time.Now().UTC()

	err := a.articles.Update(articleID, repo.ArticleUpdateInput{
		Title:     request.Title,
		Status:    request.Status,
		Icon:      request.Icon,
		Desc:      request.Desc,
		Content:   request.Content,
		UpdatedAt: UpdatedAt,
	})
	if err != nil {
		return models.Article{}, err
	}

	article, err := a.articles.GetByID(articleID)
	if err != nil {
		return models.Article{}, err
	}

	authors, err := a.articles.GetAuthors(articleID)
	if err != nil {
		return models.Article{}, err
	}
	tags, err := a.articles.GetTags(articleID)
	if err != nil {
		return models.Article{}, err
	}
	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Icon:      article.Icon,
		Desc:      article.Desc,
		Content:   article.Content,
		UpdatedAt: article.UpdatedAt,
		CreatedAt: article.CreatedAt,
		Authors:   authors,
		Tags:      tags,
	}, nil
}

func (a App) DeleteArticle(ctx context.Context, articleID uuid.UUID) error {
	return a.articles.Delete(articleID)
}

func (a App) GetArticleByID(ctx context.Context, articleID uuid.UUID) (models.Article, error) {
	article, err := a.articles.GetByID(articleID)
	if err != nil {
		return models.Article{}, err
	}

	res := models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Icon:      article.Icon,
		Desc:      article.Desc,
		Content:   article.Content,
		UpdatedAt: article.UpdatedAt,
		CreatedAt: article.CreatedAt,
	}

	authors, err := a.articles.GetAuthors(articleID)
	if err != nil {
		return models.Article{}, err
	}

	tags, err := a.articles.GetTags(articleID)
	if err != nil {
		return models.Article{}, err
	}

	res.Authors = authors
	res.Tags = tags

	return res, nil
}

func (a App) SetTags(ctx context.Context, articleID uuid.UUID, tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("too many tags")
	}
	return a.articles.SetTags(articleID, tags)
}

func (a App) SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error {
	if len(authors) > 10 {
		return fmt.Errorf("too many authors")
	}
	return a.articles.SetAuthors(articleID, authors)
}
