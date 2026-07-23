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

type ListLabelsPresenter struct {
	b.PaginatedListBuilder
	io.Writer
	e.WebPageErrorPresenter
}

func NewListLabelsPresenter(w http.ResponseWriter, p b.PageBuilder) ListLabelsPresenter {
	p.SetTitle("Labels").SetHTMLTitle("Labels").SetActiveSection(cmp.LabelsPageActive)
	pb := b.NewPaginatedListBuilder(p, listLabelsFields)
	return ListLabelsPresenter{pb, w, e.NewErrorPresenter(w)}
}

func (p ListLabelsPresenter) SuccessListLabels(r list.Response) {
	p.SetPagination(r.Pagination, rt.Labels)
	for _, l := range r.Labels {
		row := MakeListLabelRow(l)
		p.AddRow(row)
	}

	p.AddCreationButton("Create", CreateLabelForm, createLabelTargetDiv)
	p.PaginatedListBuilder.AddMarkdownPreamble(preamble)
	p.Render(p.Writer)
}

type LabelPresenter struct {
	io.Writer
	e.WebPageErrorPresenter
	successFindLabel func(lbl.Label)
}

func NewLabelPresenter(w http.ResponseWriter, mode string) LabelPresenter {
	p := LabelPresenter{Writer: w, WebPageErrorPresenter: e.NewErrorPresenter(w)}
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

func (p LabelPresenter) SuccessFindLabel(l lbl.Label) {
	p.successFindLabel(l)
}
func (p *LabelPresenter) renderEditForm(l lbl.Label) {

	b := bf.NewHTMXInlineFormBuilder(l.Name, len(listLabelsFields),
		rt.AddQueryParams(Label, "name", l.Name))
	b.AddTitle(fmt.Sprintf("Editing %v", l.Name))
	b.AddTextField("description", "Description", "description", bf.WithDefault(l.Description))
	b.Render(p.Writer)
}
func (p *LabelPresenter) renderConfirmDelete(l lbl.Label) {
	b.RenderConfirmDeleteRow(len(listLabelsFields),
		l.Name,
		"label",
		rt.AddQueryParams(Label, "name", l.Name),
		rt.AddQueryParams(Label, "name", l.Name, "mode", "view"),
		p.Writer)
}
func (p *LabelPresenter) renderView(l lbl.Label) {
	MakeListLabelRow(l).Render(p.Writer)
}

func MakeListLabelRow(l lbl.Label) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(Label, "name", l.Name, "description", l.Description, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(Label, "name", l.Name, "mode", "confirm-delete"))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(l.Name)))
	row.AddCell(tb.NewCell(Text(l.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
