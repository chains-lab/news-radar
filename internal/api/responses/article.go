package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func Article(article models.Article) resources.Article {
	var content []resources.Content
	if article.Content != nil {
		for _, c := range article.Content {
			content = append(content, resources.Content{
				Type:    string(c.Section),
				Content: c.Content,
			})
		}
	}
	data := resources.ArticleData{
		Id:   article.ID.String(),
		Type: resources.ArticleType,
		Attributes: resources.ArticleAttributes{
			Title:     article.Title,
			Status:    string(article.Status),
			Desc:      *article.Desc,
			Likes:     int64(article.Likes),
			Reposts:   int64(article.Reposts),
			CreatedAt: article.CreatedAt,
		},
	}

	if article.Content != nil {
		data.Attributes.Content = content
	}

	if article.UpdatedAt != nil {
		data.Attributes.UpdatedAt = article.UpdatedAt
	}

	if article.Authors != nil {
		var authors []resources.Relationships
		for _, author := range article.Authors {
			authors = append(authors, resources.Relationships{
				Id:   author.String(),
				Type: resources.AuthorType,
			})
		}
		data.Relationships.Authors = authors
	}

	if article.Tags != nil {
		var tags []resources.Relationships
		for _, tag := range article.Tags {
			tags = append(tags, resources.Relationships{
				Id:   tag,
				Type: resources.TagType,
			})
		}
		data.Relationships.Tags = tags
	}

	return resources.Article{
		Data: data,
	}
}
