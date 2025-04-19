package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func Author(author models.Author) resources.Author {
	res := resources.Author{
		Data: resources.AuthorData{
			Id:   author.ID.String(),
			Type: resources.AuthorType,
			Attributes: resources.AuthorAttributes{
				Name:      author.Name,
				CreatedAt: author.CreatedAt,
			},
		},
	}

	if author.UpdatedAt != nil {
		res.Data.Attributes.UpdatedAt = author.UpdatedAt
	}
	if author.Desc != nil {
		res.Data.Attributes.Desc = author.Desc
	}
	if author.Avatar != nil {
		res.Data.Attributes.Avatar = author.Avatar
	}
	if author.Twitter != nil {
		res.Data.Attributes.Twitter = author.Twitter
	}
	if author.Telegram != nil {
		res.Data.Attributes.Telegram = author.Telegram
	}
	if author.Email != nil {
		res.Data.Attributes.Email = author.Email
	}

	return res
}

func AuthorsCollection(authors []models.Author) resources.AuthorsCollection {
	data := make([]resources.AuthorData, len(authors))

	for _, author := range authors {
		element := resources.AuthorData{
			Id:   author.ID.String(),
			Type: resources.AuthorType,
			Attributes: resources.AuthorAttributes{
				Name:      author.Name,
				CreatedAt: author.CreatedAt,
			},
		}

		if author.UpdatedAt != nil {
			element.Attributes.UpdatedAt = author.UpdatedAt
		}
		if author.Desc != nil {
			element.Attributes.Desc = author.Desc
		}
		if author.Avatar != nil {
			element.Attributes.Avatar = author.Avatar
		}
		if author.Twitter != nil {
			element.Attributes.Twitter = author.Twitter
		}
		if author.Telegram != nil {
			element.Attributes.Telegram = author.Telegram
		}
		if author.Email != nil {
			element.Attributes.Email = author.Email
		}

		data = append(data, element)
	}

	return resources.AuthorsCollection{
		Data: resources.AuthorsCollectionData{
			Type: resources.AuthorsCollectionType,
			Attributes: resources.AuthorsCollectionDataAttributes{
				Data: data,
			},
		},
	}
}
