package messages

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/news-radar/internal/api/rest/responses"
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func ParseContSectionUpdateType(msg []byte) (msgType string, err error) {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(msg, &head); err != nil {
		return "", jsonkit.NewDecodeError("body", err)
	}

	return head.Type, nil
}

func ParseContentSectionUpdate(msg []byte) (req resources.UpdateContentSection, err error) {
	if err = json.Unmarshal(msg, &req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	err = validation.Errors{
		"type":    validation.Validate(req.Type, validation.Required, validation.In(resources.ContentUpdateSection)),
		"section": validation.Validate(req.Section, validation.Required),
	}.Filter()

	return req, err
}

func ParseContentSectionDelete(msg []byte) (req resources.DeleteContentSection, err error) {
	if err = json.Unmarshal(msg, &req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	err = validation.Errors{
		"type":       validation.Validate(req.Type, validation.Required, validation.In(resources.ContentDeleteSection)),
		"section_id": validation.Validate(req.SectionId, validation.Min(0)),
	}.Filter()

	return req, err
}

func ArticleContentUpdate(status string, code int, message string, section *content.Section) resources.ContentSectionResponse {
	resp := resources.ContentSectionResponse{
		Status:  status,
		Code:    int32(code),
		Message: message,
		Type:    resources.ContentUpdateSection,
	}

	if section != nil {
		resp.Section = responses.ContentSection(*section)
	}

	return resp
}
