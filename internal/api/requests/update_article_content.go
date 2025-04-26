package requests

import (
	"bytes"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleContentUpdate(msg []byte) (req resources.UpdateContentSection, err error) {
	if err = json.NewDecoder(bytes.NewReader(msg)).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return req, err
	}

	errs := validation.Errors{
		"type":    validation.Validate(req.Type, validation.Required, validation.In(resources.ContentUpdateSection)),
		"section": validation.Validate(req.Section, validation.Required),
	}
	return req, errs.Filter()
}
