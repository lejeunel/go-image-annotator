package label

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"net/http"
)

type DeleteLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(string) string
	htmx.ErrorPresenter
}

func NewDeleteLabelPresenter(w http.ResponseWriter) DeleteLabelPresenter {
	task := "Deleting label"
	okMessageFunc := func(name string) string {
		return fmt.Sprintf("Successfully deleted label %v", name)
	}
	return DeleteLabelPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteLabelPresenter) SuccessDeleteLabel(name string) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(name))
}
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.DeleteItr.Execute(r.Context(),
		r.URL.Query().Get(resourceUrlFieldName),
		NewDeleteLabelPresenter(w))
}
