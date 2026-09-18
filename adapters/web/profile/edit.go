package profile

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/update"
)

type EditProfilePresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditProfilePresenter(w http.ResponseWriter) EditProfilePresenter {
	task := "Updating profile"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated profile"
	}
	return EditProfilePresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
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
			NewDescription: r.FormValue("description"),
		},
		NewEditProfilePresenter(w))
}
