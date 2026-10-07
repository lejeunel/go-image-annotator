package role

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	se "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	"github.com/lejeunel/go-image-annotator/adapters/web/htmx"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	. "maragu.dev/gomponents"
)

var listRolesFields = []string{"name", "description", "methods", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	b.RowURL
	Writer io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PaginatedListBuilder, u b.RowURL) ListPresenter {
	return ListPresenter{p, u, w, e.NewErrorPresenter(w, p.PageBuilder)}
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

func NewViewPresenter(w http.ResponseWriter, pb b.PageBuilder, u b.RowURL) ViewPresenter {
	return ViewPresenter{Writer: w, RowURL: u, ErrorPresenter: e.NewErrorPresenter(w, pb)}
}

func (p ViewPresenter) SuccessFindRole(role r.Role) {
	MakeRow(p.RowURL, role).Render(p.Writer)
}

func (p ViewPresenter) SuccessListMethods(methods []string) {
	fmt.Println("methods", methods)
}

type DeletePresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, pb b.PageBuilder, u b.RowURL) DeletePresenter {
	return DeletePresenter{Writer: w, RowURL: u, ErrorPresenter: e.NewErrorPresenter(w, pb)}
}

func (p DeletePresenter) SuccessFindRole(role r.Role) {
	b.RenderConfirmDeleteRow(len(listRolesFields),
		role.Name, "role", p.Url, p.Writer)
}

type EditPresenter struct {
	writer http.ResponseWriter
	b.RowURL
	bf.HTMXInlineFormBuilder
	task          string
	okMessageFunc func(r.Role) string
	htmx.ErrorPresenter
	AllMethods []string
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	task := "Updating role"
	okMessageFunc := func(r r.Role) string {
		return "Successfully updated role"
	}
	return EditPresenter{
		writer: w, task: task,
		okMessageFunc:         okMessageFunc,
		HTMXInlineFormBuilder: bf.NewHTMXInlineFormBuilder(len(listRolesFields), u.Url),
		RowURL:                u,
		ErrorPresenter:        htmx.NewErrorPresenter(task, w),
	}
}

func (p *EditPresenter) SetValidMethods(methods []string) {
	p.AllMethods = methods
}

func (p *EditPresenter) SuccessFindRole(role r.Role) {
	p.HTMXInlineFormBuilder.SetResourceName(role.Name)
	p.HTMXInlineFormBuilder.AddTextField(NameFieldName, "Name", bf.WithDefault(role.Name))

	var description string
	if role.Description != nil {
		description = *role.Description
	}
	p.HTMXInlineFormBuilder.AddTextField(
		DescriptionFieldName,
		"Description",
		bf.WithDefault(description),
	)

	sb := se.NewMultiSelectBuilder(MethodsFieldName)
	for _, m := range p.AllMethods {
		sb.AddItem(m, slices.Contains(role.Methods, m))
	}

	var buf bytes.Buffer
	sb.Render(&buf)
	p.HTMXInlineFormBuilder.AddRaw("Methods", buf.String())
}

func (p EditPresenter) Render() {
	p.HTMXInlineFormBuilder.Render(p.writer)
}

func (p *EditPresenter) SuccessListMethods(methods []string) {
	sb := se.NewMultiSelectBuilder(MethodsFieldName)
	for _, m := range methods {
		sb.AddItem(m, false)
	}
	var buf bytes.Buffer
	sb.Render(&buf)
	p.AddRaw("Methods", buf.String())
}

func (p EditPresenter) SuccessUpdateRole(r r.Role) {
	htmx.NotifySuccessPayloadAndReload(p.writer, p.task, p.okMessageFunc(r))
}

func MakeRow(url b.RowURL, role r.Role) tb.Row {
	url.SetId(role.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(url.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(url.SetMode(b.ModeConfirmDelete).Url)
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(role.Name)))

	var description string
	if role.Description != nil {
		description = *role.Description
	}
	row.AddCell(tb.NewCell(Text(description)))
	row.AddCell(tb.NewCell(Text(strings.Join(role.Methods, ", "))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
