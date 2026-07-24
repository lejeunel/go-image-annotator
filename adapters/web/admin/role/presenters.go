package role

import (
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
	. "maragu.dev/gomponents"
)

var listRolesFields = []string{"name", "description", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	b.RowURL
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PaginatedListBuilder, u b.RowURL) ListPresenter {
	return ListPresenter{p, u, w, e.NewErrorPresenter(w)}
}
func (p ListPresenter) SuccessListRoles(roles []r.Role) {
	for _, role := range roles {
		row := MakeRow(p.RowURL, role)
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
func (p ViewPresenter) SuccessFindRole(role r.Role) {
	MakeRow(p.RowURL, role).Render(p.Writer)
}

type DeletePresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, u b.RowURL) DeletePresenter {
	return DeletePresenter{Writer: w, RowURL: u, ErrorPresenter: e.NewErrorPresenter(w)}
}
func (p DeletePresenter) SuccessFindRole(role r.Role) {
	b.RenderConfirmDeleteRow(len(listRolesFields),
		role.Name, "role", p.Url, p.Writer)
}

type EditPresenter struct {
	writer http.ResponseWriter
	b.RowURL
	task          string
	okMessageFunc func(update.Response) string
	htmx.ErrorPresenter
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	task := "Updating role"
	okMessageFunc := func(r update.Response) string {
		return "Successfully updated role"
	}
	return EditPresenter{writer: w, task: task,
		okMessageFunc:  okMessageFunc,
		RowURL:         u,
		ErrorPresenter: htmx.NewErrorPresenter(task, w)}
}
func (p EditPresenter) SuccessFindRole(role r.Role) {
	b := bf.NewHTMXInlineFormBuilder(role.Name, len(listRolesFields), p.Url)
	b.AddTitle(fmt.Sprintf("Editing %v", role.Name))
	b.AddTextField("name", "Name", bf.WithRequired(), bf.WithDefault(role.Name))
	b.AddTextField("description", "Description", bf.WithDefault(role.Description))
	b.Render(p.writer)

}

func (p EditPresenter) SuccessUpdateRole(r update.Response) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func MakeRow(url b.RowURL, role r.Role) tb.Row {
	url.SetId(role.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(url.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(url.SetMode(b.ModeConfirmDelete).Url)
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(role.Name)))
	row.AddCell(tb.NewCell(Text(role.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
