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
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
	. "maragu.dev/gomponents"
)

var listGroupsFields = []string{"name", "description", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PaginatedListBuilder) ListPresenter {
	return ListPresenter{p, w, e.NewErrorPresenter(w)}
}
func (p ListPresenter) SuccessListGroups(groups []g.Group) {
	for _, group := range groups {
		row := MakeRow(group)
		p.AddRow(row)
	}
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type ViewPresenter struct {
	io.Writer
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter) ViewPresenter {
	return ViewPresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p ViewPresenter) SuccessFindGroup(group g.Group) {
	MakeRow(group).Render(p.Writer)
}

type DeletePresenter struct {
	io.Writer
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter) DeletePresenter {
	return DeletePresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p DeletePresenter) SuccessFindGroup(group g.Group) {
	b.RenderConfirmDeleteRow(len(listGroupsFields),
		group.Name,
		"group",
		rt.AddQueryParams(GroupRow, "name", group.Name),
		rt.AddQueryParams(GroupRow, "name", group.Name, "mode", "view"),
		p.Writer)
}

type EditPresenter struct {
	writer        http.ResponseWriter
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditPresenter(w http.ResponseWriter) EditPresenter {
	task := "Updating group"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated group"
	}
	return EditPresenter{w, task, okMessageFunc, htmx.NewErrorPresenter(task, w)}
}
func (p EditPresenter) SuccessFindGroup(group g.Group) {
	endpoint := rt.AddQueryParams(GroupRow, "name", group.Name)
	b := bf.NewHTMXInlineFormBuilder(group.Name, len(listGroupsFields), endpoint)
	b.AddTitle(fmt.Sprintf("Editing %v", group.Name))
	b.AddTextField("name", "Name", "name", bf.WithRequired(), bf.WithDefault(group.Name))
	b.AddTextField("description", "Description", "description", bf.WithDefault(group.Description))
	b.Render(p.writer)

}

func (p EditPresenter) SuccessUpdateGroup(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func MakeRow(group g.Group) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(GroupRow, "name", group.Name, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(GroupRow, "name", group.Name,
		"mode", "confirm-delete"))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(group.Name)))
	row.AddCell(tb.NewCell(Text(group.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
