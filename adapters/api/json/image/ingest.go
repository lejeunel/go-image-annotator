package image

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	ig "github.com/lejeunel/go-image-annotator/modules/ingester"
)

type Ingest struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Ingest) Success(r ig.Response) {
	id := r.ImageId.String()
	response := models.ImageIngestionResponse{
		Id: &id,
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewIngestPresenter(w http.ResponseWriter, l slog.Logger) Ingest {
	return Ingest{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
