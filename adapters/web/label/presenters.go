package label

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	bf "github.com/lejeunel/go-image-annotator/adapters/web/builders/form"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	e "github.com/lejeunel/go-image-annotator/adapters/web/error"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	. "maragu.dev/gomponents"
)

var listLabelsFields = []string{"name", "description", "actions"}

type ListPresenter struct {
	b.PaginatedListBuilder
	b.RowURL
	io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PageBuilder, u b.RowURL) ListPresenter {
	p.SetTitle("Labels").SetHTMLTitle("Labels").SetActiveSection(cmp.LabelsPageActive)
	pb := b.NewPaginatedListBuilder(p, listLabelsFields)
	return ListPresenter{pb, u, w, e.NewErrorPresenter(w)}
}

func (p ListPresenter) SuccessListLabels(r list.Response) {
	p.SetPagination(r.Pagination, rt.Labels)
	for _, l := range r.Labels {
		row := MakeRow(p.RowURL, l)
		p.AddRow(row)
	}

	p.AddCreationButton("Create", CreateLabelFormUrl, createLabelTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type ViewPresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewViewPresenter(w http.ResponseWriter, u b.RowURL) ViewPresenter {
	return ViewPresenter{w, u, e.NewErrorPresenter(w)}
}

func (p ViewPresenter) SuccessFindLabel(l lbl.Label) {
	MakeRow(p.RowURL, l).Render(p.Writer)
}

type EditPresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewEditPresenter(w http.ResponseWriter, u b.RowURL) EditPresenter {
	return EditPresenter{w, u, e.NewErrorPresenter(w)}
}

func (p EditPresenter) SuccessFindLabel(l lbl.Label) {
	b := bf.NewHTMXInlineFormBuilder(l.Name, len(listLabelsFields), p.Url)
	b.AddTitle(fmt.Sprintf("Editing %v", l.Name))
	b.AddTextField("description", "Description", bf.WithDefault(l.Description))
	b.Render(p.Writer)
}

type DeletePresenter struct {
	io.Writer
	b.RowURL
	e.ErrorPresenter
}

func NewDeletePresenter(w http.ResponseWriter, u b.RowURL) DeletePresenter {
	return DeletePresenter{w, u, e.NewErrorPresenter(w)}
}
func (p DeletePresenter) SuccessFindLabel(l lbl.Label) {
	b.RenderConfirmDeleteRow(len(listLabelsFields),
		l.Name,
		"label",
		p.Url,
		p.Writer)
}

func MakeRow(u b.RowURL, l lbl.Label) tb.Row {
	u.SetId(l.Name)
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(u.SetMode(b.ModeEdit).Url)
	actions.SetConfirmDelete(u.SetMode(b.ModeConfirmDelete).Url)
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(l.Name)))
	row.AddCell(tb.NewCell(Text(l.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
