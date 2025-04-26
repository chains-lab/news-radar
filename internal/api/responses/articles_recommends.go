package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleRecommends(article models.Article, tags []models.Tag, authors []models.Author, recommend []models.Article) resources.ArticleWithRecommends {
	content := make([]resources.Content, 0)
	if article.Content != nil {
		for _, c := range article.Content {
			content = append(content, Content(c))
		}
	}

	data := resources.ArticleData{
		Id:   article.ID.String(),
		Type: resources.ArticleType,
		Attributes: resources.ArticleAttributes{
			Title:     article.Title,
			Status:    string(article.Status),
			CreatedAt: article.CreatedAt,
		},
	}

	if article.Content != nil {
		data.Attributes.Content = content
	}

	if article.UpdatedAt != nil {
		data.Attributes.UpdatedAt = article.UpdatedAt
	}

	if article.Desc != nil {
		data.Attributes.Desc = article.Desc
	}

	if article.Icon != nil {
		data.Attributes.Icon = article.Icon
	}

	if article.PublishedAt != nil {
		data.Attributes.PublishedAt = article.PublishedAt
	}

	authorsResp := make([]resources.AuthorData, 0)
	if authors != nil {
		for _, author := range authors {
			authorsResp = append(authorsResp, Author(author).Data)
		}
	}

	tagsResp := make([]resources.TagData, 0)
	if tags != nil {
		for _, tag := range tags {
			tagsResp = append(tagsResp, Tag(tag).Data)
		}
	}

	include := resources.ArticleWithRecommendsInclude{
		Authors: authorsResp,
		Tags:    tagsResp,
	}

	if recommend != nil {
		var res []resources.ArticleShortData
		for _, rec := range recommend {
			res = append(res, ArticleShort(rec).Data)
		}
		include.Recommends = res
	}

	return resources.ArticleWithRecommends{
		Data:     data,
		Included: include,
	}
}
