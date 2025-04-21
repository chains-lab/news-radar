package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func Tag(tag models.Tag) resources.Tag {
	return resources.Tag{
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
}

func TagsCollection(tags []models.Tag) resources.TagCollection {
	data := make([]resources.TagData, len(tags))

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
