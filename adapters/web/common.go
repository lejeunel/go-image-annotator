package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
)

type HTMXErrorPresenter struct {
	task   string
	writer http.ResponseWriter
}

func (p HTMXErrorPresenter) Error(err error) {
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

func NewHTMXErrorPresenter(t string, w http.ResponseWriter) HTMXErrorPresenter {
	return HTMXErrorPresenter{task: t, writer: w}
}

type WebPageErrorPresenter struct {
	b.PageBuilder
	writer http.ResponseWriter
}

func (p WebPageErrorPresenter) Error(err error) {
	p.PageBuilder.SetError(err).Render(p.writer)
}

func NewWebPageErrorPresenter(w http.ResponseWriter) WebPageErrorPresenter {
	return WebPageErrorPresenter{writer: w}
}
