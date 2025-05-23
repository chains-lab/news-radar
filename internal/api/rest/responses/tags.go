package responses

import (
	"github.com/chains-lab/news-radar/internal/app/models"
	"github.com/chains-lab/news-radar/resources"
)

func Tag(tag models.Tag) resources.Tag {
	res := resources.Tag{
		Data: resources.TagData{
			Id:   tag.ID,
			Type: resources.TagType,
			Attributes: resources.TagAttributes{
				Name:      tag.Name,
				Status:    string(tag.Status),
				Type:      string(tag.Type),
				Icon:      tag.Icon,
				Color:     tag.Color,
				CreatedAt: tag.CreatedAt,
			},
		},
	}
	if tag.UpdatedAt != nil {
		res.Data.Attributes.UpdatedAt = tag.UpdatedAt
	}

	return res
}

func TagsCollection(tags []models.Tag) resources.TagCollection {
	data := make([]resources.TagData, 0, len(tags))

	for _, tag := range tags {
		element := resources.TagData{
			Id:   tag.ID,
			Type: resources.TagCreateType,
			Attributes: resources.TagAttributes{
				Name:      tag.Name,
				Status:    string(tag.Status),
				Type:      string(tag.Type),
				Icon:      tag.Icon,
				Color:     tag.Color,
				CreatedAt: tag.CreatedAt,
			},
		}

		if tag.UpdatedAt != nil {
			element.Attributes.UpdatedAt = tag.UpdatedAt
		}

		data = append(data, element)
	}

	return resources.TagCollection{
		Data: resources.TagCollectionData{
			Type: resources.TagsCollectionType,
			Attributes: resources.TagCollectionDataAttributes{
				Data: data,
			},
		},
	}
}
