package web

import (
	"github.com/lejeunel/go-image-annotator/use-cases/label/update"
	"net/http"
)

type EditLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	HTMXErrorPresenter
}

func NewEditLabelPresenter(w http.ResponseWriter) EditLabelPresenter {
	task := "Updating label"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated label"
	}
	return EditLabelPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}

func (p EditLabelPresenter) SuccessUpdateLabel(r update.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) EditLabel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.Label.Update.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get("name"),
			NewDescription: r.FormValue("description"),
		},
		NewEditLabelPresenter(w))
}
