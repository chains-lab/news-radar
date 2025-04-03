package responses

import (
	"github.com/hs-zavet/news-radar/internal/app/models"
	"github.com/hs-zavet/news-radar/resources"
)

func Tag(tag models.Tag) resources.Tag {
	return resources.Tag{
		Data: resources.TagData{
			Id:   tag.Name,
			Type: resources.TagType,
			Attributes: resources.TagAttributes{
				Status:    string(tag.Status),
				Type:      string(tag.Type),
				Icon:      tag.Icon,
				Color:     tag.Color,
				CreatedAt: tag.CreatedAt,
			},
		},
	}
}
