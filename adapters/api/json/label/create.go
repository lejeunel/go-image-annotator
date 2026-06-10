package label

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
)

type Create struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Create) Success(r create.Response) {
	response := models.NewLabel{
		Name:        r.Name,
		Description: &r.Description,
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewCreatePresenter(w http.ResponseWriter, l slog.Logger) Create {
	return Create{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
