package collection

import (
	"fmt"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
)

type CreateCollectionPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(create.Response) string
	profiles      []string
	htmx.ErrorPresenter
}

func NewCreateCollectionPresenter(w http.ResponseWriter) CreateCollectionPresenter {
	task := "Creating collection"
	okMessageFunc := func(r create.Response) string {
		return fmt.Sprintf("Successfully created collection %v", r.Name)
	}
	return CreateCollectionPresenter{
		writer: w, task: task,
		okMessageFunc: okMessageFunc, ErrorPresenter: htmx.NewErrorPresenter(task, w),
	}
}

func (p CreateCollectionPresenter) SuccessCreateCollection(r create.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (p *CreateCollectionPresenter) SuccessListAllProfiles(profiles []string) {
	p.profiles = profiles
}

func (p CreateCollectionPresenter) Render() {
	b := bf.NewHTMXCreateFormBuilder(CollectionUrl, createCollectionTargetDiv)
	b.AddTitle("Create a new collection")
	b.AddTextField(nameFieldName, "Name", bf.WithRequired())
	b.AddTextField(descriptionFieldName, "Description")

	profilePicker := se.NewSingleSelect(p.profiles,
		profileFieldName)
	b.AddRaw(profileFieldLabel, profilePicker)
	b.Render(p.writer)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	s.CreateItr.Execute(r.Context(),
		create.Request{
			Name:        r.FormValue(nameFieldName),
			Description: r.FormValue(descriptionFieldName),
		}, NewCreateCollectionPresenter(w))
}

func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	p := NewCreateCollectionPresenter(w)
	s.ListProfileItr.Execute(r.Context(), &p)
	p.Render()
}
