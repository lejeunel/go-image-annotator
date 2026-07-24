package group

import (
	"fmt"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"net/http"
)

var createGroupTargetDiv = "create-group"

type CreateGroupPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	htmx.ErrorPresenter
}

func NewCreateGroupPresenter(w http.ResponseWriter) CreateGroupPresenter {
	task := "Creating group"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created group %v", r.Name)
	}
	return CreateGroupPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}
func (p CreateGroupPresenter) Success(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Groups.Create.Execute(r.Context(), create.Request{Name: r.FormValue("name"), Description: r.FormValue("description")},
		NewCreateGroupPresenter(w))
}
