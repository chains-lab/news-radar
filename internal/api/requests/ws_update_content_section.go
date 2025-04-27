package requests

import (
	"encoding/json"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/news-radar/resources"
)

func ParseArticleContentWS(msg []byte) (msgType string, payload any, err error) {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(msg, &head); err != nil {
		return "", nil, jsonkit.NewDecodeError("body", err)
	}

	switch head.Type {
	case resources.ContentUpdateSection:
		var req resources.UpdateContentSection
		if err := json.Unmarshal(msg, &req); err != nil {
			return head.Type, nil, jsonkit.NewDecodeError("body", err)
		}

		err = validation.Errors{
			"type":    validation.Validate(req.Type, validation.Required, validation.In(resources.ContentUpdateSection)),
			"section": validation.Validate(req.Section, validation.Required),
		}.Filter()

		return head.Type, req, err

	case resources.ContentDeleteSection:
		var req resources.DeleteContentSection
		if err := json.Unmarshal(msg, &req); err != nil {
			return head.Type, nil, jsonkit.NewDecodeError("body", err)
		}

		err = validation.Errors{
			"type":       validation.Validate(req.Type, validation.Required, validation.In(resources.ContentDeleteSection)),
			"section_id": validation.Validate(req.SectionId, validation.Min(0)),
		}.Filter()

		return head.Type, req, err

	default:
		return head.Type, nil, fmt.Errorf("unsupported ws message type %q", head.Type)
	}
}
