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
	pg "github.com/lejeunel/go-image-annotator/adapters/web/pagination"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	rt "github.com/lejeunel/go-image-annotator/routes"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	. "maragu.dev/gomponents"
)

var listLabelsFields = []string{"name", "description", "actions"}

//go:embed preamble.md
var preamble string

type ListLabelsPresenter struct {
	b.PaginatedListBuilder
	io.Writer
	e.WebPageErrorPresenter
}

func NewListLabelsPresenter(w http.ResponseWriter, p b.PageBuilder) ListLabelsPresenter {
	p.SetTitle("Labels").SetHTMLTitle("Labels").SetActiveSection(cmp.LabelsPageActive)
	b := b.NewPaginatedListBuilder(p, listLabelsFields)
	return ListLabelsPresenter{b, w, e.NewErrorPresenter(w)}
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
func (p ListLabelsPresenter) SuccessFindLabel(l l.Label) {
	MakeListLabelRow(l).Render(p.Writer)
}
func (s *Server) TableRow(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	switch r.URL.Query().Get("mode") {
	case "confirm-delete":
		b.RenderConfirmDeleteRow(len(listLabelsFields),
			name,
			"label",
			rt.AddQueryParams(Label, "name", name),
			rt.AddQueryParams(Label, "name", name, "mode", "view"),
			w)
	case "edit":
		b := bf.NewHTMXInlineFormBuilder(name, len(listLabelsFields),
			rt.AddQueryParams(Label, "name", name))
		b.AddTitle(fmt.Sprintf("Editing %v", name))
		b.AddTextField("description", "Description", "description", bf.WithDefault(r.URL.Query().Get("description")))
		b.Render(w)
	default:
		s.FindItr.Execute(r.Context(),
			name,
			NewListLabelsPresenter(w, s.PageBuilder))
	}
}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.ListItr.Execute(r.Context(),
		pag.PaginationParams{PageSize: s.DefaultPageSize, Page: pg.GetPageFromRequest(r)},
		NewListLabelsPresenter(w, s.PageBuilder))
}

func MakeListLabelRow(l l.Label) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(Label, "name", l.Name, "description", l.Description, "mode", "edit"))
	actions.SetConfirmDelete(rt.AddQueryParams(Label, "name", l.Name, "mode", "confirm-delete"))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(l.Name)))
	row.AddCell(tb.NewCell(Text(l.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
