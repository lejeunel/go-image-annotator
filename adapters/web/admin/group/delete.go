package group

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
)

type DeleteGroupPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(string) string
	htmx.ErrorPresenter
}

func NewDeleteGroupPresenter(w http.ResponseWriter) DeleteGroupPresenter {
	task := "deleting group"
	okMessageFunc := func(name string) string {
		return fmt.Sprintf("Successfully deleted group %v", name)
	}
	return DeleteGroupPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteGroupPresenter) SuccessDeleteGroup(name string) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(name))
}
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.Groups.Delete.Execute(r.Context(),
		r.URL.Query().Get(resourceUrlFieldName),
		NewDeleteGroupPresenter(w))
}
