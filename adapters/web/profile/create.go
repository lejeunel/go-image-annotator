package profile

import (
	"bytes"
	"fmt"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/create"
)

type CreateProfilePresenter struct {
	writer          http.ResponseWriter
	task            string
	okMessageFunc   func(pr.Profile) string
	availableLabels []string
	htmx.ErrorPresenter
}

func NewCreateProfilePresenter(w http.ResponseWriter) CreateProfilePresenter {
	task := "Creating profile"
	okMessageFunc := func(p pr.Profile) string {
		return fmt.Sprintf("Successfully created profile %v", p.Name)
	}
	return CreateProfilePresenter{
		writer: w,
		task:   task, okMessageFunc: okMessageFunc,
		ErrorPresenter: htmx.NewErrorPresenter(task, w),
	}
}

func (p CreateProfilePresenter) SuccessCreateProfile(profile pr.Profile) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(profile))
}

func (p *CreateProfilePresenter) SuccessFetchLabels(labels []string) {
	p.availableLabels = labels
}

func (p *CreateProfilePresenter) Render() {
	b := bf.NewHTMXCreateFormBuilder(ProfileUrl, createProfileTargetDiv)
	b.AddTitle("Create a new profile")
	b.AddTextField(NameFieldName, "Name", bf.WithRequired())
	b.AddTextField(DescriptionFieldName, "Description")

	sb := se.NewMultiSelectBuilder(LabelsFieldName)
	for _, l := range p.availableLabels {
		sb.AddItem(l, false)
	}
	var buf bytes.Buffer
	sb.Render(&buf)
	b.AddRaw("Labels", buf.String())
	b.Render(p.writer)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.CreateItr.Execute(r.Context(), create.Request{
		Name:        r.FormValue(NameFieldName),
		Description: r.FormValue(DescriptionFieldName),
		Labels:      r.Form[LabelsFieldName],
	}, NewCreateProfilePresenter(w))
}

func (s *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	p := NewCreateProfilePresenter(w)
	s.ListAllLabelsItr.Execute(r.Context(), &p)
	p.Render()
}
