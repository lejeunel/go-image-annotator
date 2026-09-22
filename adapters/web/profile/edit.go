package profile

import (
	"net/http"

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
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
	b.RowURL
	availableLabels []string
}

func NewEditProfilePresenter(w http.ResponseWriter, url b.RowURL) EditProfilePresenter {
	task := "Updating profile"
	okMessageFunc := func(r update.Response) string {
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

func (p EditProfilePresenter) SuccessFindProfile(l pr.Profile) {
	b := bf.NewHTMXInlineFormBuilder(len(listProfilesFields), p.Url)
	labelPicker := se.NewMultiSelectCombobox(p.availableLabels,
		LabelsFieldName)
	b.AddRaw("Labels", labelPicker)

	b.SetResourceName(l.Name)
	b.AddTextField("description", "Description", bf.WithDefault(l.Description))
	b.Render(p.writer)
}

func (p EditProfilePresenter) SuccessUpdateProfile(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.UpdateItr.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get(resourceUrlFieldName),
			NewDescription: r.FormValue(DescriptionFieldName),
		},
		NewEditProfilePresenter(w, s.RowURL))
}
