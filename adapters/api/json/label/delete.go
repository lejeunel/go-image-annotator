package label

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
)

type Delete struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Delete) SuccessDeleteLabel(string) {
	p.Writer.WriteHeader(http.StatusNoContent)

}

func NewDeletePresenter(w http.ResponseWriter, l slog.Logger) Delete {
	return Delete{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
