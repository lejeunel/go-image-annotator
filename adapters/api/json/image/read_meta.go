package image

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type ReadMeta struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p ReadMeta) SuccessReadImage(image im.Image) {
	response := BuildImageResponse(image)
	json.WriteJSON(p.Writer, 200, response)

}

func NewReadMetaPresenter(w http.ResponseWriter, l slog.Logger) ReadMeta {
	return ReadMeta{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
