package collection

import (
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func (s *Server) Edit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}
	req := update.Request{
		Name:           r.URL.Query().Get(resourceUrlFieldName),
		NewName:        r.FormValue(nameFieldName),
		NewDescription: r.FormValue(descriptionFieldName),
	}

	group := r.FormValue(groupFieldName)
	if group != publicGroupPlaceholderName {
		req.NewGroup = &group
	}

	s.UpdateItr.Execute(r.Context(), req,
		NewEditPresenter(w, s.RowURL))
}

type EditPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	Form          bf.HTMXInlineFormBuilder
	htmx.ErrorPresenter
	groupOfCollection *string
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	task := "Updating collection"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated collection"
	}
	form := bf.NewHTMXInlineFormBuilder(len(listCollectionsFields), u.Url)
	return EditPresenter{w, task, okMessageFunc, form, htmx.NewErrorPresenter(task, w), nil}
}

func (p EditPresenter) SuccessUpdateCollection(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func (p *EditPresenter) SuccessFindCollection(c clc.Collection) {
	if c.Group != nil {
		p.groupOfCollection = c.Group
	}
	p.Form.SetResourceName(c.Name)
	p.Form.AddTextField("name", "Name", bf.WithRequired(), bf.WithDefault(c.Name))
	p.Form.AddTextField("description", "Description", bf.WithDefault(c.Description))
}

func (p *EditPresenter) SuccessListGroups(groups []grp.Group) {
	cb := p.Form.AddCombobox("Group", groupFieldName)
	cb.AddField(publicGroupPlaceholderName)
	for _, group := range groups {
		cb.AddField(group.Name)
	}
	if p.groupOfCollection != nil {
		cb.SetSelectedValue(*p.groupOfCollection)
	}
	p.Form.Render(p.writer)
}
