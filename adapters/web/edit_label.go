package web

import (
	"fmt"
	"net/http"

	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/label/update"
)

type EditLabelPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	HTMXErrorPresenter
}

func NewEditLabelPresenter(w http.ResponseWriter) EditLabelPresenter {
	task := "Updating label"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated label"
	}
	return EditLabelPresenter{w, task, okMessageFunc, NewHTMXErrorPresenter(task, w)}
}

func (p EditLabelPresenter) SuccessUpdateLabel(r update.Response) {
	payload, _ := NotifySuccessPayload(p.task, p.okMessageFunc(r))
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusOK)
}
func (s *Server) EditLabelForm(w http.ResponseWriter, r *http.Request) {
	currentName := r.URL.Query().Get("name")
	endpoint := rt.AddQueryParams(rt.Label, "name", currentName)
	b := bf.NewHTMXInlineFormBuilder(len(listLabelsFields), endpoint, bf.HTMXPutMethod)
	b.AddTitle(fmt.Sprintf("Editing %v", currentName))
	b.AddTextField("description", "Description", "description", bf.WithDefault(r.URL.Query().Get("description")))
	b.Render(w)
}
func (s *Server) EditLabel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.Label.Update.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get("name"),
			NewDescription: r.FormValue("description"),
		},
		NewEditLabelPresenter(w))
}
