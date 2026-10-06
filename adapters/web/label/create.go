package label

import (
	"fmt"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/use-cases/label/create"
)

type CreateLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(l.Label) string
	htmx.ErrorPresenter
}

func NewCreateLabelPresenter(w http.ResponseWriter) CreateLabelPresenter {
	task := "Creating label"
	okMessageFunc := func(r l.Label) string {
		return fmt.Sprintf("Successfully created label %v", r.Name)
	}
	return CreateLabelPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p CreateLabelPresenter) Success(r l.Label) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	var description *string
	descriptionValue := r.FormValue(createDescriptionFieldName)
	if descriptionValue != "" {
		description = &descriptionValue
	}
	s.CreateItr.Execute(r.Context(), create.Request{
		Name:        r.FormValue(createNameFieldName),
		Description: description,
	}, NewCreateLabelPresenter(w))
}

func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXFormBuilder(LabelUrl, createLabelTargetDiv)
	b.AddTitle("Create a new label")
	b.AddTextField(createNameFieldName, "Name", bf.WithRequired())
	b.AddTextField(createDescriptionFieldName, "Description")
	b.Render(w)
}
