package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/gatekit/jsonkit"
	"github.com/chains-lab/news-radar/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateArticle(r *http.Request) (req resources.CreateArticle, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = jsonkit.NewDecodeError("body", err)
		return req, err
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.ArticleCreateType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
