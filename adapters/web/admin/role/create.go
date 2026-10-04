package role

import (
	"bytes"
	"fmt"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/role/create"
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

func (p CreateRolePresenter) SuccessCreateRole(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Roles.Create.Execute(r.Context(),
		create.Request{
			Name:        r.FormValue(NameFieldName),
			Description: r.FormValue(DescriptionFieldName),
			Methods:     r.Form[MethodsFieldName],
		},
		NewCreateRolePresenter(w))
}

func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXFormBuilder(RoleRowUrl, createRoleTargetDiv)
	b.AddTitle("Create a new role")
	b.AddTextField(NameFieldName, "Name", bf.WithRequired())
	b.AddTextField(DescriptionFieldName, "Description")

	sb := se.NewMultiSelectBuilder(MethodsFieldName)
	for _, m := range s.Auth.ListMethods() {
		sb.AddItem(m, false)
	}
	var buf bytes.Buffer
	sb.Render(&buf)
	b.AddRaw("Methods", buf.String())
	b.Render(w)
}
