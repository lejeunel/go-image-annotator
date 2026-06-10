package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
)

type List struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p List) Success(r list.Response) {
	data := []models.Collection{}
	for _, c := range r.Collections {
		data = append(data,
			models.Collection{
				Name:        &c.Name,
				Description: &c.Description,
			})
	}

	response := models.ListCollectionsResponse{Data: &data,
		Pagination: json.BuildPaginationResponse(r.Pagination),
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewListPresenter(w http.ResponseWriter, l slog.Logger) List {
	return List{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
