package user

import (
	"net/http"
	"strings"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

type EditUserPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditUserPresenter(w http.ResponseWriter) EditUserPresenter {
	task := "Updating user"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated user"
	}
	return EditUserPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditUserPresenter) SuccessUpdate(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	var groups []string
	groupsStr := r.FormValue(groupsFieldName)
	if groupsStr != "" {
		groups = strings.Split(groupsStr, ",")
	}

	var roles []string
	rolesStr := r.FormValue(rolesFieldName)
	if rolesStr != "" {
		roles = strings.Split(rolesStr, ",")
	}

	s.Users.UpdatePrivileges.Execute(r.Context(),
		update.Request{
			Id:     r.URL.Query().Get(resourceUrlFieldName),
			Groups: groups,
			Roles:  roles,
		},
		NewEditUserPresenter(w))
}
