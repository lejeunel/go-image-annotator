package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
)

type Presenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Presenter) SuccessFindCollection(r clc.Collection) {
	response := models.Collection{
		Name:        r.Name,
		Description: r.Description,
		Group:       r.Group,
		Profile:     r.Profile,
	}

	json.WriteJSON(p.Writer, 200, response)
}

func (p Presenter) SuccessSubmitCloneTask(r clone.Response) {
	response := models.TaskResponse{
		TaskId: r.Id, Issuer: r.Issuer, Type: r.Type,
	}

	json.WriteJSON(p.Writer, 200, response)
}

func NewPresenter(w http.ResponseWriter, l slog.Logger) Presenter {
	return Presenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
