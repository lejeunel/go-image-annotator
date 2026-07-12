package web

import (
	"fmt"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
	"net/http"
)

var createLabelTargetDiv = "create-label"

type CreateLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	HTMXErrorPresenter
}

func NewCreateLabelPresenter(w http.ResponseWriter) CreateLabelPresenter {
	task := "Creating label"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created label %v", r.Name)
	}
	return CreateLabelPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}
func (p CreateLabelPresenter) Success(r create.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) CreateLabel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.Label.Create.Execute(r.Context(), create.Request{Name: r.FormValue("name"),
		Description: r.FormValue("description")}, NewCreateLabelPresenter(w))
}
func (s *Server) CreateLabelForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(rt.Label, createLabelTargetDiv)
	b.AddTitle("Create a new label")
	b.AddTextField("name", "Name", "name", bf.WithRequired())
	b.AddTextField("description", "Description", "description")
	b.Render(w)
}
