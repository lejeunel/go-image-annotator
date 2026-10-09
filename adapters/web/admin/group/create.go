package group

import (
	"fmt"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
)

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

func (p CreateGroupPresenter) SuccessCreateGroup(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	req := create.Request{
		Name: r.FormValue(NameFieldName),
	}

	d := r.FormValue(DescriptionFieldName)
	if d != "" {
		req.Description = &d
	}
	s.Groups.Create.Execute(r.Context(),
		req,
		NewCreateGroupPresenter(w))
}
