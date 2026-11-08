package htmx

import (
	"fmt"
	"net/http"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ErrorPresenter struct {
	task   string
	writer http.ResponseWriter
}

func (p ErrorPresenter) Error(err error) {
	NotifyError(p.writer, fmt.Sprintf("failed %v", p.task), e.PublicMessage(err))
}

func NewErrorPresenter(t string, w http.ResponseWriter) ErrorPresenter {
	return ErrorPresenter{task: t, writer: w}
}
