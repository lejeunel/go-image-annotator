package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
)

type Find struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Find) SuccessFindCollection(r clc.Collection) {
	response := models.Collection{
		Name:        r.Name,
		Description: &r.Description,
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewFindPresenter(w http.ResponseWriter, l slog.Logger) Find {
	return Find{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
