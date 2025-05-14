package app

import (
	"context"
	"errors"

	"github.com/chains-lab/news-radar/internal/app/ape"
	"github.com/chains-lab/news-radar/internal/app/models"
	"github.com/chains-lab/news-radar/internal/content"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a App) UpdateArticleContent(ctx context.Context, articleID uuid.UUID, sections []content.Section) (models.Article, error) {
	article, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.Article{}, ape.ErrArticleNotFound
		default:
			return models.Article{}, err
		}
	}

	for idx, sec := range sections {
		if sec.ID != idx {
			return models.Article{}, ape.ErrArticleContentNumerationIsIncorrect
		}
	}

	for i, _ := range article.Content {
		err = a.articles.DeleteContentSection(articleID, i)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				return models.Article{}, ape.ErrArticleNotFound
			default:
				return models.Article{}, err
			}
		}
	}

	for _, section := range sections {
		if section.Text == nil && section.Media == nil && section.Audio == nil {
			return models.Article{}, errors.New("no content section")
		} else {
			err := a.articles.UpdateContentSection(articleID, section)
			if err != nil {
				switch {
				case errors.Is(err, mongo.ErrNoDocuments):
					return models.Article{}, ape.ErrArticleNotFound
				default:
					return models.Article{}, err
				}
			}
		}
	}

	return a.GetArticleByID(ctx, articleID)
}

func (a App) UpdateContentSection(ctx context.Context, articleID uuid.UUID, section content.Section) (models.Article, error) {
	if section.Text == nil && section.Media == nil && section.Audio == nil {
		return models.Article{}, errors.New("no content section")
	} else {
		err := a.articles.UpdateContentSection(articleID, section)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				return models.Article{}, ape.ErrArticleNotFound
			default:
				return models.Article{}, err
			}
		}
	}

	article, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.Article{}, ape.ErrArticleNotFound
		default:
			return models.Article{}, err
		}
	}

	res := ArticleRepoToModels(article)
	return res, nil
}

func (a App) DeleteContentSection(ctx context.Context, articleID uuid.UUID, sectionID int) (models.Article, error) {
	err := a.articles.DeleteContentSection(articleID, sectionID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.Article{}, ape.ErrArticleNotFound
		default:
			return models.Article{}, err
		}
	}

	article, err := a.articles.GetByID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.Article{}, ape.ErrArticleNotFound
		default:
			return models.Article{}, err
		}
	}

	res := ArticleRepoToModels(article)
	return res, nil
}
