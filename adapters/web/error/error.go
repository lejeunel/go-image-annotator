package error

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"net/http"
)

type ErrorPresenter struct {
	b.PageBuilder
	writer http.ResponseWriter
}

func (p ErrorPresenter) Error(err error) {
	p.PageBuilder.SetError(err).Render(p.writer)
}

func NewErrorPresenter(w http.ResponseWriter) ErrorPresenter {
	return ErrorPresenter{writer: w}
}
