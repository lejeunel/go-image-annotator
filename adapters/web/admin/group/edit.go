package group

import (
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

type EditGroupPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditGroupPresenter(w http.ResponseWriter) EditGroupPresenter {
	task := "Updating group"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated group"
	}
	return EditGroupPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditGroupPresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	req := update.Request{
		Name:    r.URL.Query().Get(resourceUrlFieldName),
		NewName: r.FormValue("name"),
	}
	d := r.FormValue(DescriptionFieldName)
	if d != "" {
		req.NewDescription = &d
	}
	s.Groups.Update.Execute(r.Context(),
		req,
		NewEditPresenter(w, s.RowUrl))
}
