package label

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
)

type List struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p List) SuccessListLabels(r list.Response) {
	labels := []models.Label{}
	for _, label := range r.Labels {
		labels = append(labels, MakeLabelResponse(label))
	}

	s.WriteJSON(p.Writer, 200, models.ListLabels{
		Labels:     labels,
		Pagination: json.BuildPaginationResponse(r.Pagination),
	})
}

func NewListPresenter(w http.ResponseWriter, l slog.Logger) List {
	return List{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
