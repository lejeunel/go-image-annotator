package json

import (
	"log/slog"
	"net/http"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type MetaPresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func (p MetaPresenter) SuccessAddMetadata(im im.Image) {
	response := BuildImageResponse(im)
	s.WriteJSON(p.Writer, 200, response)
}

func (p MetaPresenter) SuccessDeleteMetadata(key string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func NewMetaPresenter(w http.ResponseWriter, l slog.Logger) MetaPresenter {
	return MetaPresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
