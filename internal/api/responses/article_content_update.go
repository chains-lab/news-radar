package responses

import (
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleContentUpdate(status string, message string, section *content.Section) resources.ContentUpdateResponse {
	resp := resources.ContentUpdateResponse{
		Data: resources.ContentUpdateResponseData{
			Type: resources.ArticleContentUpdateResponseType,
			Attributes: resources.ContentUpdateResponseDataAttributes{
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
