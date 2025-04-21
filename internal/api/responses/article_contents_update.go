package responses

import (
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleContentUpdate(status string, message string, section *content.Section) resources.UpdateContentResponse {
	resp := resources.UpdateContentResponse{
		Data: resources.UpdateContentResponseData{
			Type: resources.ArticleContentUpdateResponseType,
			Attributes: resources.UpdateContentResponseDataAttributes{
				Status:  status,
				Message: message,
			},
		},
	}
	
	if section != nil {
		ctnt := Content(*section)
		resp.Data.Attributes.Content = &ctnt
	}

	return resp
}
