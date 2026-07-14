package error

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"net/http"
)

type WebPageErrorPresenter struct {
	b.PageBuilder
	writer http.ResponseWriter
}

func (p WebPageErrorPresenter) Error(err error) {
	p.PageBuilder.SetError(err).Render(p.writer)
}

func NewErrorPresenter(w http.ResponseWriter) WebPageErrorPresenter {
	return WebPageErrorPresenter{writer: w}
}
