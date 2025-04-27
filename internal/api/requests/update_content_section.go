package requests

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/jsonkit"
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
