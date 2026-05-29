package label

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"net/http"
)

type Delete struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Delete) Success() {
	p.Writer.WriteHeader(http.StatusNoContent)

}

func NewDeletePresenter(w http.ResponseWriter) Delete {
	return Delete{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
