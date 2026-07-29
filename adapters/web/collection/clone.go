package collection

import (
	"fmt"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
)

func (s *Server) Clone(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	var withAnnotations bool
	if r.FormValue(cloneWithAnnotationsFieldName) == "on" {
		withAnnotations = true
	}

	source := r.URL.Query().Get(resourceUrlFieldName)
	s.CloneItr.Execute(r.Context(),
		clone.Request{
			Source:      source,
			Destination: r.FormValue(cloneNameFieldName),
			Deep:        withAnnotations,
		},
		NewClonePresenter(w, s.RowURL))

	s.FindItr.Execute(r.Context(), source,
		NewViewPresenter(w, s.RowURL))
}

type ClonePresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(clone.Response) string
	htmx.ErrorPresenter
}

func NewClonePresenter(w http.ResponseWriter, u b.RowURL) ClonePresenter {
	task := "Cloning collection"
	okMessageFunc := func(r clone.Response) string {
		return fmt.Sprintf("Started cloning task with id %v", r.Id)
	}
	return ClonePresenter{w, u, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}

func (p ClonePresenter) SuccessFindCollection(c clc.Collection) {
	b := bf.NewHTMXInlineFormBuilder(len(listCollectionsFields), p.Url, bf.WithMode(bf.CloneMode))
	b.SetResourceName(c.Name)
	b.AddTextField(cloneNameFieldName, "Name", bf.WithRequired(), bf.WithDefault(c.Name))
	b.AddTextField(cloneDescriptionFieldName, "Description", bf.WithDefault(c.Description))
	b.AddCheckbox(cloneWithAnnotationsFieldName, "With annotations")
	b.Render(p.writer)
}

func (p ClonePresenter) SuccessSubmitCloneTask(r clone.Response) {
	htmx.NotifySuccessPayload(p.writer, p.task, p.okMessageFunc(r))
}
