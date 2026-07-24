package role

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type DeleteRolePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(u.UserId) string
	htmx.ErrorPresenter
}

func NewDeleteRolePresenter(w http.ResponseWriter) DeleteRolePresenter {
	task := "deleting role"
	okMessageFunc := func(id u.UserId) string {
		return fmt.Sprintf("Successfully deleted role %v", id)
	}
	return DeleteRolePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteRolePresenter) SuccessDeleteRole(name string) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(name))
}
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.Roles.Delete.Execute(r.Context(),
		r.URL.Query().Get(resourceUrlFieldName),
		NewDeleteRolePresenter(w))
}
