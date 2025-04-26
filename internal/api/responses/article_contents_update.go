package responses

import (
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleContentUpdate(status string, message string, section *content.Section) resources.UpdateContentSectionResponse {
	resp := resources.UpdateContentSectionResponse{
		Status:  status,
		Message: message,
		Type:    resources.ContentUpdateSection,
	}

	if section != nil {
		resp.Section = Content(*section)
	}

	return resp
}
