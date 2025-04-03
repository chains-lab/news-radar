package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/config"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/internal/repo"
)

type Articles struct {
	data articlesRepo
}

type articlesRepo interface {
	Create(input repo.ArticleCreateInput) error
	Update(ID uuid.UUID, input repo.ArticleUpdateInput) error
	Delete(ID uuid.UUID) error

	SetTags(ID uuid.UUID, tags []string) error
	AddTag(ID uuid.UUID, tag string) error
	DeleteTag(ID uuid.UUID, tag string) error

	AddAuthor(ID uuid.UUID, author uuid.UUID) error
	DeleteAuthor(ID uuid.UUID, author uuid.UUID) error
	SetAuthors(ID uuid.UUID, authors []uuid.UUID) error

	GetByID(ID uuid.UUID) (repo.ArticleModel, error)
}

func NewArticles(cfg config.Config) (*Articles, error) {
	data, err := repo.NewArticles(cfg)
	if err != nil {
		return nil, err
	}

	return &Articles{
		data: data,
	}, nil
}

type CreateArticleRequest struct {
	ID      uuid.UUID         `json:"id"`
	Title   string            `json:"title"`
	Icon    string            `json:"icon"`
	Desc    string            `json:"desc"`
	Content []content.Section `json:"content,omitempty"`
}

func (a *Articles) CreateArticle(ctx context.Context, request CreateArticleRequest) error {
	ArticleID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.data.Create(repo.ArticleCreateInput{
		ID:        ArticleID,
		Title:     request.Title,
		Icon:      request.Icon,
		Desc:      request.Desc,
		Content:   request.Content,
		Status:    string(models.AuthorStatusInactive),
		CreatedAt: CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

type UpdateArticleRequest struct {
	Title    *string           `json:"title,omitempty"`
	Icon     *string           `json:"icon,omitempty"`
	Desc     *string           `json:"desc,omitempty"`
	Content  []content.Section `json:"content,omitempty"`
	Status   *string           `json:"status,omitempty"`
	Likes    *int              `json:"likes,omitempty"`
	Reposts  *int              `json:"reposts,omitempty"`
	Dislikes *int              `json:"dislikes,omitempty"`
}

func (a *Articles) UpdateArticle(ctx context.Context, articleID uuid.UUID, request UpdateArticleRequest) error {
	UpdatedAt := time.Now().UTC()

	_, err := models.ParseArticleStatus(*request.Status)
	if err != nil {
		return err
	}
	return a.data.Update(articleID, repo.ArticleUpdateInput{
		Title:     request.Title,
		Icon:      request.Icon,
		Desc:      request.Desc,
		Content:   request.Content,
		Status:    request.Status,
		Dislike:   request.Dislikes,
		Likes:     request.Likes,
		Reposts:   request.Reposts,
		UpdatedAt: UpdatedAt,
	})
}

func (a *Articles) DeleteArticle(ctx context.Context, articleID uuid.UUID) error {
	return a.data.Delete(articleID)
}

func (a *Articles) GetByID(ctx context.Context, articleID uuid.UUID) (models.Article, error) {
	article, err := a.data.GetByID(articleID)
	if err != nil {
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Icon:      article.Icon,
		Desc:      article.Desc,
		Content:   article.Content,
		Likes:     article.Likes,
		Reposts:   article.Reposts,
		Status:    models.ArticleStatus(article.Status),
		UpdatedAt: article.UpdatedAt,
		CreatedAt: article.CreatedAt,
	}, nil
}

func (a *Articles) SetTags(ctx context.Context, articleID uuid.UUID, tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("too many tags")
	}
	return a.data.SetTags(articleID, tags)
}

func (a *Articles) AddTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	return a.data.AddTag(articleID, tag)
}

func (a *Articles) DeleteTag(ctx context.Context, articleID uuid.UUID, tag string) error {
	return a.data.DeleteTag(articleID, tag)
}

func (a *Articles) AddAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	return a.data.AddAuthor(articleID, authorID)
}

func (a *Articles) DeleteAuthor(ctx context.Context, articleID uuid.UUID, authorID uuid.UUID) error {
	return a.data.DeleteAuthor(articleID, authorID)
}
