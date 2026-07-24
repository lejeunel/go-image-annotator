package user

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"net/http"
)

type CreateUserPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	htmx.ErrorPresenter
}

func NewCreateUserPresenter(w http.ResponseWriter) CreateUserPresenter {
	task := "Creating user"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created user %v", r.Id)
	}
	return CreateUserPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}
func (p CreateUserPresenter) Success(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	var isAdmin bool
	if r.FormValue(createIsAdminFieldName) == "true" {
		isAdmin = true
	}
	s.Users.Create.Execute(r.Context(),
		create.Request{Id: r.FormValue(createEmailFieldName), IsAdmin: isAdmin},
		NewCreateUserPresenter(w))
}
