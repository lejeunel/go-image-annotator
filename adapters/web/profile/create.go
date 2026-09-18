package profile

import (
	"fmt"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/create"
)

type CreateProfilePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(pr.Profile) string
	htmx.ErrorPresenter
}

func NewCreateProfilePresenter(w http.ResponseWriter) CreateProfilePresenter {
	task := "Creating profile"
	okMessageFunc := func(p pr.Profile) string {
		return fmt.Sprintf("Successfully created profile %v", p.Name)
	}
	return CreateProfilePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p CreateProfilePresenter) SuccessCreateProfile(profile pr.Profile) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(profile))
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.CreateItr.Execute(r.Context(), create.Request{
		Name:        r.FormValue(createNameFieldName),
		Description: r.FormValue(createDescriptionFieldName),
		Labels:      r.Form[createLabelsFieldName],
	}, NewCreateProfilePresenter(w))
}

func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	b := bf.NewHTMXCreateFormBuilder(ProfileUrl, createProfileTargetDiv)
	b.AddTitle("Create a new profile")
	b.AddTextField(createNameFieldName, "Name", bf.WithRequired())
	b.AddTextField(createDescriptionFieldName, "Description")
	labelPicker := se.NewMultiLabelCombobox([]string{"first-label", "second-label"},
		createLabelsFieldName)
	b.AddRaw("Labels", labelPicker)
	b.Render(w)
}
