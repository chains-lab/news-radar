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
	Name string `json:"name" bson:"name"`
}

func (a App) CreateAuthor(ctx context.Context, request CreateAuthorRequest) (models.Author, error) {
	AuthorID := uuid.New()
	CreatedAt := time.Now().UTC()

	author, err := a.authors.Create(repo.AuthorCreateInput{
		ID:        AuthorID,
		Name:      request.Name,
		Status:    enums.AuthorStatusActive,
		CreatedAt: CreatedAt,
	})
	if err != nil {
		return models.Author{}, err
	}

	return models.Author{
		ID:        author.ID,
		Name:      author.Name,
		Status:    author.Status,
		CreatedAt: author.CreatedAt,
	}, nil
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

func (a App) UpdateAuthor(ctx context.Context, authorID uuid.UUID, request UpdateAuthorRequest) (models.Author, error) {
	var author repo.AuthorModel

	input := repo.AuthorUpdateInput{}
	updated := false

	if request.Name != nil {
		input.Name = request.Name
		updated = true
	}

	if request.Status != nil {
		input.Status = request.Status
		updated = true
	}

	if request.Desc != nil {
		input.Desc = request.Desc
		updated = true
	}

	if request.Avatar != nil {
		input.Avatar = request.Avatar
		updated = true
	}

	if request.Email != nil || request.Telegram != nil || request.Twitter != nil {
		input.Email = request.Email
		updated = true
	}

	if !updated {
		return a.GetAuthorByID(ctx, authorID)
	}

	return models.Author{
		ID:        author.ID,
		Name:      author.Name,
		Status:    author.Status,
		Desc:      author.Desc,
		Avatar:    author.Avatar,
		Email:     author.Email,
		Telegram:  author.Telegram,
		Twitter:   author.Twitter,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
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

func (a App) GetArticlesForAuthor(ctx context.Context, articleID uuid.UUID) ([]models.Article, error) {
	articles, err := a.articles.GetArticlesForAuthor(articleID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Article, 0, len(articles))

	for _, article := range articles {
		thisArticle, err := a.articles.GetByID(article)
		if err != nil {
			return nil, err
		}

		elem := models.Article{
			ID:        thisArticle.ID,
			Title:     thisArticle.Title,
			Status:    thisArticle.Status,
			CreatedAt: thisArticle.CreatedAt,
		}

		if thisArticle.Desc != nil {
			elem.Desc = thisArticle.Desc
		}

		if thisArticle.Icon != nil {
			elem.Icon = thisArticle.Icon
		}

		if thisArticle.Content != nil {
			elem.Content = thisArticle.Content
		}

		if thisArticle.UpdatedAt != nil {
			elem.UpdatedAt = thisArticle.UpdatedAt
		}

		authors, err := a.articles.GetAuthors(article)
		if err != nil {
			return nil, err
		}

		tags, err := a.articles.GetTags(article)
		if err != nil {
			return nil, err
		}

		elem.Authors = authors
		elem.Tags = tags

		res = append(res, elem)
	}

	return res, nil
}
