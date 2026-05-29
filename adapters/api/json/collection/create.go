package collection

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
)

type Create struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Create) Success(r create.Response) {
	response := models.NewCollection{
		Name:        r.Name,
		Description: &r.Description,
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewCreatePresenter(w http.ResponseWriter) Create {
	return Create{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
