package collection

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

type Update struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p Update) Success(r update.Response) {
	p.Writer.WriteHeader(http.StatusOK)
}

func NewUpdatePresenter(w http.ResponseWriter) Update {
	return Update{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
