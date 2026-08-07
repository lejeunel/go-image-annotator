package label

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	l "github.com/lejeunel/go-image-annotator/entities/label"
)

type Find struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Find) SuccessFindLabel(r l.Label) {
	response := models.Label{
		Name:        &r.Name,
		Description: &r.Description,
	}

	json.WriteJSON(p.Writer, 200, response)
}

func NewFindPresenter(w http.ResponseWriter, l slog.Logger) Find {
	return Find{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
