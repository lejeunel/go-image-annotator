package role

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/role/create"
	"net/http"
)

type CreateRolePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	htmx.ErrorPresenter
}

func NewCreateRolePresenter(w http.ResponseWriter) CreateRolePresenter {
	task := "Creating role"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created role %v", r.Name)
	}
	return CreateRolePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}
func (p CreateRolePresenter) Success(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Roles.Create.Execute(r.Context(),
		create.Request{
			Name:        r.FormValue(createNameFieldName),
			Description: r.FormValue(createDescriptionFieldName)},
		NewCreateRolePresenter(w))
}
