package web

import (
	"fmt"
	"net/http"

	rt "github.com/lejeunel/go-image-annotator/routes"
)

type DeleteLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(string) string
	HTMXErrorPresenter
}

func NewDeleteLabelPresenter(w http.ResponseWriter) DeleteLabelPresenter {
	task := "Deleting label"
	okMessageFunc := func(name string) string {
		return fmt.Sprintf("Successfully deleted label %v", name)
	}
	return DeleteLabelPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}

func (p DeleteLabelPresenter) SuccessDeleteLabel(name string) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(name))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) ConfirmDeleteLabel(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	RenderConfirmDeleteRow(len(listLabelsFields),
		name,
		"label",
		rt.AddQueryParams(rt.Label, "name", name),
		w)
}
func (s *Server) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	s.Label.Delete.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewDeleteLabelPresenter(w))
}
