package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleShort(article models.Article) resources.ArticleShort {
	res := resources.ArticleShort{
		Data: resources.ArticleShortData{
			Type: resources.ArticleShortCollectionType,
			Id:   article.ID.String(),
			Attributes: resources.ArticleShortDataAttributes{
				Title:     article.Title,
				CreatedAt: article.CreatedAt,
			},
		},
	}
	if article.Icon != nil {
		res.Data.Attributes.Icon = *article.Icon
	} else {
		res.Data.Attributes.Icon = ""
	}

	if article.Desc != nil {
		res.Data.Attributes.Desc = *article.Desc
	} else {
		res.Data.Attributes.Desc = ""
	}

	return res
}

func ArticleShortsCollection(articles []models.Article) resources.ArticleShortCollection {
	data := make([]resources.ArticleShortData, 0, len(articles))

	for _, article := range articles {
		element := resources.ArticleShortData{
			Type: resources.ArticleShortCollectionType,
			Id:   article.ID.String(),
			Attributes: resources.ArticleShortDataAttributes{
				Title:     article.Title,
				CreatedAt: article.CreatedAt,
			},
		}
		if article.Icon != nil {
			element.Attributes.Icon = *article.Icon
		} else {
			element.Attributes.Icon = ""
		}

		if article.Desc != nil {
			element.Attributes.Desc = *article.Desc
		} else {
			element.Attributes.Desc = ""
		}

		data = append(data, element)
	}

	return resources.ArticleShortCollection{
		Data: resources.ArticleShortCollectionData{
			Type: resources.ArticleShortCollectionType,
			Attributes: resources.ArticleShortCollectionDataAttributes{
				Data: data,
			},
		},
	}
}
