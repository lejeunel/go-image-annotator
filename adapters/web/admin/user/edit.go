package user

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

type EditUserPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditUserPresenter(w http.ResponseWriter) EditUserPresenter {
	task := "Updating collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated collection"
	}
	return EditUserPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditUserPresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}
func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	// s.Users.UpdatePrivileges.Execute(r.Context(),
	// 	update.Request{
	// 		Id:             r.URL.Query().Get("id"),
	// 		Groups:        r.FormValue("name"),
	// 		NewDescription: r.FormValue("description"),
	// 	},
	// 	NewEditUserPresenter(w))
}
