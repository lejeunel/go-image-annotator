package profile

import (
	"bytes"
	"net/http"
	"slices"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/update"
)

type EditProfilePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(pr.Profile) string
	htmx.ErrorPresenter
	b.RowURL
	availableLabels []string
}

func NewEditProfilePresenter(w http.ResponseWriter, url b.RowURL) EditProfilePresenter {
	task := "Updating profile"
	okMessageFunc := func(r pr.Profile) string {
		return "Successfully updated profile"
	}
	return EditProfilePresenter{
		writer: w,
		task:   task, okMessageFunc: okMessageFunc, ErrorPresenter: htmx.NewErrorPresenter(task, w),
		RowURL: url,
	}
}

func (p *EditProfilePresenter) SuccessFetchLabels(labels []string) {
	p.availableLabels = labels
}

func (p EditProfilePresenter) SuccessFindProfile(profile pr.Profile) {
	b := bf.NewHTMXInlineFormBuilder(len(listProfilesFields), p.Url)
	sb := se.NewMultiSelectBuilder(LabelsFieldName)
	for _, l := range p.availableLabels {
		sb.AddItem(l, slices.Contains(profile.Labels, l))
	}
	var buf bytes.Buffer
	sb.Render(&buf)
	b.AddRaw("Labels", buf.String())
	b.SetResourceName(profile.Name)
	b.AddTextField(NameFieldName, "Name", bf.WithDefault(profile.Name))

	var description string
	if profile.Description != nil {
		description = *profile.Description
	}
	b.AddTextField(DescriptionFieldName, "Description", bf.WithDefault(description))
	b.Render(p.writer)
}

func (p EditProfilePresenter) SuccessUpdateProfile(r pr.Profile) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	var description *string
	descriptionField := r.FormValue(DescriptionFieldName)
	if descriptionField != "" {
		description = &descriptionField
	}
	s.UpdateItr.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get(resourceUrlFieldName),
			NewName:        r.FormValue(NameFieldName),
			NewDescription: description,
			NewLabels:      r.Form[LabelsFieldName],
		},
		NewEditProfilePresenter(w, s.RowURL))
}
