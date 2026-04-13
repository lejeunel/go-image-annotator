package image

import (
	"github.com/lejeunel/go-image-annotator-v2/adapters/api/json"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	"net/http"
)

type ReadMeta struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p ReadMeta) SuccessReadImage(image im.Image) {
	response := BuildImageResponse(image)
	json.WriteJSON(p.Writer, 200, response)

}

func NewReadMetaPresenter(w http.ResponseWriter) ReadMeta {
	return ReadMeta{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
