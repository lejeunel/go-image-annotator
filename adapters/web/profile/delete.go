package profile

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
)

type DeleteProfilePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(string) string
	htmx.ErrorPresenter
}

func NewDeleteProfilePresenter(w http.ResponseWriter) DeleteProfilePresenter {
	task := "Deleting profile"
	okMessageFunc := func(name string) string {
		return fmt.Sprintf("Successfully deleted profile %v", name)
	}
	return DeleteProfilePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteProfilePresenter) SuccessDeleteProfile(name string) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(name))
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.DeleteItr.Execute(r.Context(),
		r.URL.Query().Get(resourceUrlFieldName),
		NewDeleteProfilePresenter(w))
}
