package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
)

type Presenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func MakeCollectionResponse(c clc.Collection) models.Collection {
	return models.Collection{
		Name:        c.Name,
		Description: c.Description,
		Group:       c.Group,
		Profile:     c.Profile,
	}
}

func (p Presenter) SuccessFindCollection(r clc.Collection) {
	s.WriteJSON(p.Writer, 200, MakeCollectionResponse(r))
}

func (p Presenter) SuccessSubmitCloneTask(r clone.Response) {
	response := models.TaskResponse{
		TaskId: r.Id, Issuer: r.Issuer, Type: r.Type,
	}

	s.WriteJSON(p.Writer, 200, response)
}

func NewPresenter(w http.ResponseWriter, l slog.Logger) Presenter {
	return Presenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
