package htmx

import (
	"fmt"
	"net/http"
)

type ErrorPresenter struct {
	task   string
	writer http.ResponseWriter
}

func (p ErrorPresenter) Error(err error) {
	NotifyError(p.writer, fmt.Sprintf("failed %v", p.task), err.Error())
}

func NewErrorPresenter(t string, w http.ResponseWriter) ErrorPresenter {
	return ErrorPresenter{task: t, writer: w}
}
