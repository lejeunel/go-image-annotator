package role

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
)

type EditRolePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditRolePresenter(w http.ResponseWriter) EditRolePresenter {
	task := "Updating role"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated role"
	}
	return EditRolePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditRolePresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Roles.Update.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get(resourceUrlFieldName),
			NewName:        r.FormValue("name"),
			NewDescription: r.FormValue("description"),
		},
		NewEditPresenter(w, s.RowUrl))
}
