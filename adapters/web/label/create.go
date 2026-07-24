package label

import (
	"fmt"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
	"net/http"
)

type CreateLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	htmx.ErrorPresenter
}

func NewCreateLabelPresenter(w http.ResponseWriter) CreateLabelPresenter {
	task := "Creating label"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created label %v", r.Name)
	}
	return CreateLabelPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}
func (p CreateLabelPresenter) Success(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.CreateItr.Execute(r.Context(), create.Request{Name: r.FormValue("name"),
		Description: r.FormValue("description")}, NewCreateLabelPresenter(w))
}
func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(LabelUrl, createLabelTargetDiv)
	b.AddTitle("Create a new label")
	b.AddTextField(createNameFieldName, "Name", bf.WithRequired())
	b.AddTextField(createDescriptionFieldName, "Description")
	b.Render(w)
}
