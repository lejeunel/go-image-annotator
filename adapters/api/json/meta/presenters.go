package meta

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	jim "github.com/lejeunel/go-image-annotator/adapters/api/json/image"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	s "github.com/lejeunel/go-image-annotator/shared"
)

type Presenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Presenter) SuccessAddMetadata(im im.Image) {
	response := jim.BuildImageResponse(im)
	s.WriteJSON(p.Writer, 200, response)
}

func (p Presenter) SuccessDeleteMetadata(key string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func NewPresenter(w http.ResponseWriter, l slog.Logger) Presenter {
	return Presenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
