package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/comtools/jsonkit"
	"github.com/hs-zavet/news-radar/resources"
)

func UpdateAuthor(r *http.Request) (req resources.UpdateAuthor, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.AuthorUpdateType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
