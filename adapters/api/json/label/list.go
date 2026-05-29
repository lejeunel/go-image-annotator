package label

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"net/http"
)

type List struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p List) Success(r list.Response) {
	data := []models.Label{}
	for _, label := range r.Labels {
		data = append(data,
			models.Label{
				Name:        &label.Name,
				Description: &label.Description,
			})
	}

	response := models.ListLabelsResponse{Data: &data,
		Pagination: json.BuildPaginationResponse(r.Pagination),
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewListPresenter(w http.ResponseWriter) List {
	return List{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
