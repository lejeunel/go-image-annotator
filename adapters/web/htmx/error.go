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
	payload, _ := NotifyError(fmt.Sprintf("failed %v", p.task),
		err.Error())
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusUnprocessableEntity)
}

func NewErrorPresenter(t string, w http.ResponseWriter) ErrorPresenter {
	return ErrorPresenter{task: t, writer: w}
}
