package collection

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/read"
	"net/http"
)

type Find struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Find) Success(r read.Response) {
	response := models.Collection{
		Name:        &r.Name,
		Description: &r.Description,
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewFindPresenter(w http.ResponseWriter) Find {
	return Find{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
