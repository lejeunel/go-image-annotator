package collection

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
)

type Delete struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Delete) Success(delete.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)

}

func NewDeletePresenter(w http.ResponseWriter, l slog.Logger) Delete {
	return Delete{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
