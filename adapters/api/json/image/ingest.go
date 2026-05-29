package image

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"net/http"
)

type Ingest struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Ingest) Success(r ingest.Response) {
	id := r.ImageId.String()
	response := models.ImageIngestionResponse{
		Id: &id,
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewIngestPresenter(w http.ResponseWriter) Ingest {
	return Ingest{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
