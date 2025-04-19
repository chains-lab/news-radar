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
	Title  *string              `json:"title,omitempty"`
	Status *enums.ArticleStatus `json:"status,omitempty"`
	Icon   *string              `json:"icon,omitempty"`
	Desc   *string              `json:"desc,omitempty"`
}

func (a App) UpdateArticle(ctx context.Context, articleID uuid.UUID, request UpdateArticleRequest) (models.Article, error) {
	var article repo.ArticleModel
	var err error

	updated := false

	if request.Title != nil {
		article, err = a.articles.UpdateTitle(articleID, *request.Title)
		if err != nil {
			return models.Article{}, err
		}
		updated = true
	}

	if request.Status != nil {
		article, err = a.articles.UpdateStatus(articleID, *request.Status)
		if err != nil {
			return models.Article{}, err
		}
		updated = true
	}

	if request.Icon != nil {
		article, err = a.articles.UpdateIcon(articleID, *request.Icon)
		if err != nil {
			return models.Article{}, err
		}
		updated = true
	}

	if request.Desc != nil {
		article, err = a.articles.UpdateDesc(articleID, *request.Desc)
		if err != nil {
			return models.Article{}, err
		}
		updated = true
	}

	if !updated {
		//for idempotency
		return a.GetArticleByID(ctx, articleID)
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

func (a App) UpdateArticleContent(ctx context.Context, articleID uuid.UUID, index int, section content.Section) (models.Article, error) {
	res, err := a.articles.UpdateContent(articleID, index, section)
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
		ID:        res.ID,
		Title:     res.Title,
		Icon:      res.Icon,
		Desc:      res.Desc,
		Content:   res.Content,
		UpdatedAt: res.UpdatedAt,
		CreatedAt: res.CreatedAt,
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

//HASHTAGS

func (a App) SetArticleTags(ctx context.Context, articleID uuid.UUID, tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("too many tags")
	}
	return a.articles.SetTags(articleID, tags)
}

func (a App) GetArticleForTags(ctx context.Context, tag string) ([]models.Article, error) {
	articles, err := a.articles.GetArticlesForTag(tag)
	if err != nil {
		return nil, err
	}

	var res []models.Article
	for _, articleID := range articles {
		article, err := a.articles.GetByID(articleID)
		if err != nil {
			return nil, err
		}

		elem := models.Article{
			ID:        article.ID,
			Status:    article.Status,
			Title:     article.Title,
			CreatedAt: article.CreatedAt,
		}

		if article.Desc != nil {
			elem.Desc = article.Desc
		}

		if article.Icon != nil {
			elem.Icon = article.Icon
		}

		if article.Content != nil {
			elem.Content = article.Content
		}

		if article.UpdatedAt != nil {
			elem.UpdatedAt = article.UpdatedAt
		}

		authors, err := a.articles.GetAuthors(articleID)
		if err != nil {
			return nil, err
		}

		tags, err := a.articles.GetTags(articleID)
		if err != nil {
			return nil, err
		}

		elem.Authors = authors
		elem.Tags = tags

		res = append(res, elem)
	}

	return res, nil
}

func (a App) GetArticleTags(ctx context.Context, articleID uuid.UUID) ([]models.Tag, error) {
	tags, err := a.articles.GetTags(articleID)
	if err != nil {
		return nil, err
	}

	var res []models.Tag
	for _, tagID := range tags {
		tag, err := a.tags.Get(tagID)
		if err != nil {
			return nil, err
		}

		res = append(res, models.Tag{
			Name:      tag.Name,
			Status:    tag.Status,
			Type:      tag.Type,
			Color:     tag.Color,
			Icon:      tag.Icon,
			CreatedAt: tag.CreatedAt,
		})
	}

	return res, nil
}

func (a App) AddArticleTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	err := a.articles.AddTag(articleID, tag)
	if err != nil {
		return err
	}

	return nil
}

func (a App) DeleteArticleTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	err := a.articles.DeleteTag(articleID, tag)
	if err != nil {
		return err
	}

	return nil
}

func (a App) CleanArticleTags(ctx context.Context, articleID uuid.UUID) error {
	tags, err := a.articles.GetTags(articleID)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		err := a.articles.DeleteTag(articleID, tag)
		if err != nil {
			return err
		}
	}

	return nil
}

//AUTHORSHIP

func (a App) SetAuthors(ctx context.Context, articleID uuid.UUID, authors []uuid.UUID) error {
	if len(authors) > 10 {
		return fmt.Errorf("too many authors")
	}
	return a.articles.SetAuthors(articleID, authors)
}

func (a App) GetArticleForAuthors(ctx context.Context, authorID uuid.UUID) ([]models.Article, error) {
	articles, err := a.articles.GetArticlesForAuthor(authorID)
	if err != nil {
		return nil, err
	}

	var res []models.Article
	for _, articleID := range articles {
		article, err := a.articles.GetByID(articleID)
		if err != nil {
			return nil, err
		}

		elem := models.Article{
			ID:        article.ID,
			Status:    article.Status,
			Title:     article.Title,
			CreatedAt: article.CreatedAt,
		}

		if article.Desc != nil {
			elem.Desc = article.Desc
		}

		if article.Icon != nil {
			elem.Icon = article.Icon
		}

		if article.Content != nil {
			elem.Content = article.Content
		}

		if article.UpdatedAt != nil {
			elem.UpdatedAt = article.UpdatedAt
		}

		authors, err := a.articles.GetAuthors(articleID)
		if err != nil {
			return nil, err
		}

		tags, err := a.articles.GetTags(articleID)
		if err != nil {
			return nil, err
		}

		elem.Authors = authors
		elem.Tags = tags

		res = append(res, elem)
	}

	return res, nil
}

func (a App) AddArticleAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	err := a.articles.AddAuthor(articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) DeleteArticleAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	err := a.articles.DeleteAuthor(articleID, authorID)
	if err != nil {
		return err
	}

	return nil
}

func (a App) CleanArticleAuthors(ctx context.Context, articleID uuid.UUID) error {
	authors, err := a.articles.GetArticlesForAuthor(articleID)
	if err != nil {
		return err
	}

	for _, author := range authors {
		err := a.articles.DeleteAuthor(articleID, author)
		if err != nil {
			return err
		}
	}

	return nil
}
