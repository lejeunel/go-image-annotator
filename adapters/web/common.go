package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebErrorPresenter struct {
	task   string
	writer http.ResponseWriter
}

func (p WebErrorPresenter) Error(err error) {
	payload, _ := json.Marshal(map[string]any{
		"htmx-notify": map[string]string{
			"variant": "danger",
			"title":   fmt.Sprintf("failed %v", p.task),
			"message": err.Error(),
		},
	})
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusUnprocessableEntity)
}

func NewWebErrorPresenter(t string, w http.ResponseWriter) WebErrorPresenter {
	return WebErrorPresenter{task: t, writer: w}
}
