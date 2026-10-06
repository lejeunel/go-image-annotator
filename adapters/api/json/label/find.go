package label

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Find struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func MakeLabelResponse(l l.Label) models.Label {
	return models.Label{
		Name:        l.Name,
		Description: l.Description,
	}
}

func (p Find) SuccessFindLabel(r l.Label) {
	s.WriteJSON(p.Writer, 200, MakeLabelResponse(r))
}

func NewFindPresenter(w http.ResponseWriter, l slog.Logger) Find {
	return Find{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
