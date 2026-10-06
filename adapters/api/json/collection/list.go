package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
)

type List struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p List) SuccessListCollections(r list.Response) {
	collections := []models.Collection{}
	for _, c := range r.Collections {
		collections = append(collections,
			MakeCollectionResponse(c))
	}

	response := models.ListCollections{
		Collections: collections,
		Pagination:  json.BuildPaginationResponse(r.Pagination),
	}

	s.WriteJSON(p.Writer, 200, response)
}

func NewListPresenter(w http.ResponseWriter, l slog.Logger) List {
	return List{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
