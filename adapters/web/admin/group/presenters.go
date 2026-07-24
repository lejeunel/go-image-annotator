package group

import (
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
	. "maragu.dev/gomponents"
)

var listGroupsFields = []string{"name", "description", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	b.RowURL
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PaginatedListBuilder, u b.RowURL) ListPresenter {
	return ListPresenter{p, u, w, e.NewErrorPresenter(w)}
}
func (p ListPresenter) SuccessListGroups(groups []g.Group) {
	for _, group := range groups {
		row := MakeRow(p.RowURL, group)
		p.AddRow(row)
	}
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type ViewPresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, u b.RowURL) ViewPresenter {
	return ViewPresenter{Writer: w, RowURL: u, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p ViewPresenter) SuccessFindGroup(group g.Group) {
	MakeRow(p.RowURL, group).Render(p.Writer)
}

type DeletePresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, u b.RowURL) DeletePresenter {
	return DeletePresenter{Writer: w, RowURL: u, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p DeletePresenter) SuccessFindGroup(group g.Group) {
	b.RenderConfirmDeleteRow(len(listGroupsFields),
		group.Name, "group", p.Url, p.Writer)
}

type EditPresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	task := "Updating group"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated group"
	}
	return EditPresenter{writer: w, task: task,
		okMessageFunc:  okMessageFunc,
		RowURL:         u,
		ErrorPresenter: htmx.NewErrorPresenter(task, w)}
}
func (p EditPresenter) SuccessFindGroup(group g.Group) {
	b := bf.NewHTMXInlineFormBuilder(group.Name, len(listGroupsFields), p.Url)
	b.AddTitle(fmt.Sprintf("Editing %v", group.Name))
	b.AddTextField("name", "Name", bf.WithRequired(), bf.WithDefault(group.Name))
	b.AddTextField("description", "Description", bf.WithDefault(group.Description))
	b.Render(p.writer)

}

func (p EditPresenter) SuccessUpdateGroup(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func MakeRow(url b.RowURL, group g.Group) tb.Row {
	url.SetId(group.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(url.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(url.SetMode(b.ModeConfirmDelete).Url)
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(group.Name)))
	row.AddCell(tb.NewCell(Text(group.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
