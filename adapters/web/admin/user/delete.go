package user

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type DeleteUserPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(u.UserId) string
	htmx.ErrorPresenter
}

func NewDeleteUserPresenter(w http.ResponseWriter) DeleteUserPresenter {
	task := "deleting user"
	okMessageFunc := func(id u.UserId) string {
		return fmt.Sprintf("Successfully deleted user %v", id)
	}
	return DeleteUserPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p DeleteUserPresenter) SuccessDeleteUser(id u.UserId) {
	payload, _ := htmx.NotifySuccessPayloadAndReload(p.task, p.okMessageFunc(id))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	s.Users.Delete.Execute(r.Context(),
		r.URL.Query().Get("id"),
		NewDeleteUserPresenter(w))
}
