package json

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
)

type LabelPresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func MakeLabelResponse(l l.Label) models.Label {
	return models.Label{
		Name:        l.Name,
		Description: l.Description,
	}
}

func (p LabelPresenter) SuccessFindLabel(r l.Label) {
	s.WriteJSON(p.Writer, 200, MakeLabelResponse(r))
}

func (p LabelPresenter) SuccessCreateLabel(r l.Label) {
	s.WriteJSON(p.Writer, 200, MakeLabelResponse(r))
}

func (p LabelPresenter) SuccessDeleteLabel(string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p LabelPresenter) SuccessListLabels(r list.Response) {
	labels := []models.Label{}
	for _, label := range r.Labels {
		labels = append(labels, MakeLabelResponse(label))
	}

	s.WriteJSON(p.Writer, 200, models.ListLabels{
		Labels:     labels,
		Pagination: BuildPaginationResponse(r.Pagination),
	})
}

func NewLabelPresenter(w http.ResponseWriter, l slog.Logger) LabelPresenter {
	return LabelPresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
