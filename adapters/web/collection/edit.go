package collection

import (
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	s.UpdateItr.Execute(r.Context(),
		update.Request{
			Name:           r.URL.Query().Get(resourceUrlFieldName),
			NewName:        r.FormValue(createNameFieldName),
			NewDescription: r.FormValue(createDescriptionFieldName),
		},
		NewEditPresenter(w, s.RowURL))
}

type EditPresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	task := "Updating collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated collection"
	}
	return EditPresenter{w, u, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p EditPresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (p EditPresenter) SuccessFindCollection(c clc.Collection) {
	b := bf.NewHTMXInlineFormBuilder(c.Name, len(listCollectionsFields), p.Url)
	b.AddTextField("name", "Name", bf.WithRequired(), bf.WithDefault(c.Name))
	b.AddTextField("description", "Description", bf.WithDefault(c.Description))
	b.Render(p.writer)
}
