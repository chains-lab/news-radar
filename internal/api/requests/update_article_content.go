package requests

import (
	"bytes"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/news-radar/resources"
)

func ArticleContentUpdate(msg []byte) (req resources.UpdateContent, err error) {
	if err = json.NewDecoder(bytes.NewReader(msg)).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return req, err
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.ArticleContentUpdateType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
