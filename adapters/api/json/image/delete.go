package image

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
)

type Delete struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Delete) SuccessDeleteImage(r delete.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func NewDeleteImagePresenter(w http.ResponseWriter, l slog.Logger) Delete {
	return Delete{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
