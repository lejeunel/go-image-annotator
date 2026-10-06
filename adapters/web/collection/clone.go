package collection

import (
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	s "github.com/lejeunel/go-image-annotator/adapters/web/shared"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
)

func (s *Server) Clone(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	var deep bool
	if r.FormValue(deepFieldName) == "on" {
		deep = true
	}

	source := r.URL.Query().Get(resourceUrlFieldName)
	s.CloneItr.Execute(r.Context(),
		clone.Request{
			Source:      source,
			Destination: r.FormValue(nameFieldName),
			Deep:        &deep,
		},
		NewClonePresenter(w, s.RowURL))

	s.FindItr.Execute(r.Context(), source,
		NewViewPresenter(w, s.PageBuilder, s.RowURL))
}

type ClonePresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(clone.Response) string
	htmx.ErrorPresenter
	groups []string
}

func NewClonePresenter(w http.ResponseWriter, u b.RowURL) ClonePresenter {
	task := "Cloning collection"
	okMessageFunc := func(r clone.Response) string {
		return s.MakeNewTaskMessage()
	}
	return ClonePresenter{
		writer: w, RowURL: u, task: task,
		okMessageFunc: okMessageFunc, ErrorPresenter: htmx.NewErrorPresenter(task, w),
	}
}

func (p ClonePresenter) SuccessSubmitCloneTask(r clone.Response) {
	htmx.NotifySuccessPayload(p.writer, p.task, p.okMessageFunc(r))
}

func (p *ClonePresenter) SuccessListGroups(groups []g.Group) {
	for _, g := range groups {
		p.groups = append(p.groups, g.Name)
	}
}

func (p ClonePresenter) SuccessFindCollection(c clc.Collection) {
	b := bf.NewHTMXInlineFormBuilder(len(listCollectionsFields), p.Url, bf.WithMode(bf.CloneMode))
	b.SetResourceName(c.Name)
	b.AddTextField(nameFieldName, "Name", bf.WithRequired(), bf.WithDefault(c.Name))

	if c.Description != nil {
		b.AddTextField(descriptionFieldName, "Description", bf.WithDefault(*c.Description))
	} else {
		b.AddTextField(descriptionFieldName, "Description")
	}

	groupSelect := se.NewSingleSelect(p.groups, groupFieldName)
	b.AddRaw(groupFieldName, groupSelect)
	b.AddCheckbox(deepFieldName, "Deep")
	b.Render(p.writer)
}
