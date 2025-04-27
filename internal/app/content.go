package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hs-zavet/news-radar/internal/app/ape"
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/internal/content"
	"go.mongodb.org/mongo-driver/mongo"
)

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
