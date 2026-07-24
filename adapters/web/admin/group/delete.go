package group

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type DeleteGroupPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(u.UserId) string
	htmx.ErrorPresenter
}

func NewDeleteGroupPresenter(w http.ResponseWriter) DeleteGroupPresenter {
	task := "deleting group"
	okMessageFunc := func(id u.UserId) string {
		return fmt.Sprintf("Successfully deleted group %v", id)
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
