package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func Article(article models.Article) resources.Article {
	content := make([]resources.Content, 0)
	if article.Content != nil {
		for _, c := range article.Content {
			section := resources.Content{
				Id:   c.ID,
				Type: string(c.Type),
			}

			if c.Media != nil {
				section.Media = &resources.ContentMedia{
					Url:     c.Media.URL,
					Caption: c.Media.Caption,
					Alt:     c.Media.Alt,
					Width:   int32(c.Media.Width),
					Height:  int32(c.Media.Height),
					Source:  c.Media.Source,
				}
			}

			if c.Text != nil {
				text := make([]resources.ContentTextInner, 0)
				for _, t := range c.Text {
					if t.Text != nil {
						marks := make([]string, 0)
						for _, mark := range t.Marks {
							marks = append(marks, string(mark))
						}

						text = append(text, resources.ContentTextInner{
							Text:  t.Text,
							Color: t.Color,
							Link:  t.Link,
							Marks: marks,
						})
					}
				}
				section.Text = text
			}

			content = append(content, section)
		}
	}

	data := resources.ArticleData{
		Id:   article.ID.String(),
		Type: resources.ArticleType,
		Attributes: resources.ArticleAttributes{
			Title:     article.Title,
			Status:    string(article.Status),
			Desc:      *article.Desc,
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
