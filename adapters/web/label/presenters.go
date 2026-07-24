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
	io.Writer
	e.ErrorPresenter
}

func NewListPresenter(w http.ResponseWriter, p b.PageBuilder) ListPresenter {
	p.SetTitle("Labels").SetHTMLTitle("Labels").SetActiveSection(cmp.LabelsPageActive)
	pb := b.NewPaginatedListBuilder(p, listLabelsFields)
	return ListPresenter{pb, w, e.NewErrorPresenter(w)}
}

func (p ListPresenter) SuccessListLabels(r list.Response) {
	p.SetPagination(r.Pagination, rt.Labels)
	for _, l := range r.Labels {
		row := MakeRow(l)
		p.AddRow(row)
	}

	p.AddCreationButton("Create", CreateLabelForm, createLabelTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type RowPresenter struct {
	io.Writer
	e.ErrorPresenter
	successFindLabel func(lbl.Label)
}

func NewLabelPresenter(w http.ResponseWriter, mode string) RowPresenter {
	p := RowPresenter{Writer: w, ErrorPresenter: e.NewErrorPresenter(w)}
	switch mode {
	case "edit":
		p.successFindLabel = p.renderEditForm
	case "confirm-delete":
		p.successFindLabel = p.renderConfirmDelete
	default:
		p.successFindLabel = p.renderView
	}
	return p
}

func (p RowPresenter) SuccessFindLabel(l lbl.Label) {
	p.successFindLabel(l)
}
func (p *RowPresenter) renderEditForm(l lbl.Label) {

	b := bf.NewHTMXInlineFormBuilder(l.Name, len(listLabelsFields),
		rt.AddQueryParams(Label, "name", l.Name))
	b.AddTitle(fmt.Sprintf("Editing %v", l.Name))
	b.AddTextField("description", "Description", "description", bf.WithDefault(l.Description))
	b.Render(p.Writer)
}
func (p *RowPresenter) renderConfirmDelete(l lbl.Label) {
	b.RenderConfirmDeleteRow(len(listLabelsFields),
		l.Name,
		"label",
		rt.AddQueryParams(Label, "name", l.Name),
		rt.AddQueryParams(Label, "name", l.Name, "mode", "view"),
		p.Writer)
}
func (p *RowPresenter) renderView(l lbl.Label) {
	MakeRow(l).Render(p.Writer)
}

func MakeRow(l lbl.Label) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(Label, "name", l.Name, "description", l.Description, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(Label, "name", l.Name, "mode", "confirm-delete"))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(l.Name)))
	row.AddCell(tb.NewCell(Text(l.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
