package responses

import (
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleContentUpdate(status string, code int, message string, section *content.Section) resources.ContentSectionResponse {
	resp := resources.ContentSectionResponse{
		Status:  status,
		Code:    int32(code),
		Message: message,
		Type:    resources.ContentUpdateSection,
	}

	if section != nil {
		resp.Section = ContentSection(*section)
	}

	return resp
}
