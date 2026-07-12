package web

import (
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	. "maragu.dev/gomponents"
)

var listLabelsFields = []string{"name", "description", "actions"}

type ListLabelsPresenter struct {
	b.PaginatedListBuilder
	io.Writer
	WebPageErrorPresenter
}

func NewListLabelsPresenter(w http.ResponseWriter, p b.PageBuilder) ListLabelsPresenter {
	p.SetTitle("Labels").SetActive(cmp.LabelsPageActive)
	b := b.NewPaginatedListBuilder(p, listLabelsFields)
	return ListLabelsPresenter{b, w, NewWebPageErrorPresenter(w)}
}

func (p ListLabelsPresenter) SuccessListLabels(r list.Response) {
	p.SetPagination(r.Pagination, rt.Labels)
	for _, l := range r.Labels {
		row := MakeListLabelRow(l)
		p.AddRow(row)
	}

	p.AddCreationButton("Create new label", rt.CreateLabelForm, createLabelTargetDiv)
	p.Build().Render(p.Writer)
}
func (p ListLabelsPresenter) SuccessFindLabel(l l.Label) {
	MakeListLabelRow(l).Render(p.Writer)
}

func (s *Server) GetLabel(w http.ResponseWriter, r *http.Request) {
	s.Label.Find.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewListLabelsPresenter(w, s.PageBuilder))
}

func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.Label.List.Execute(r.Context(),
		list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListLabelsPresenter(w, s.PageBuilder))
}

func MakeListLabelRow(l l.Label) tb.Row {
	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(rt.EditLabelForm, "name", l.Name))
	actions.SetConfirmDelete(rt.AddQueryParams(rt.ConfirmDeleteLabel, "name", l.Name))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(Text(l.Name)))
	row.AddCell(tb.NewCell(Text(l.Description)))
	row.AddCell(tb.NewCell(actions.Build()))
	return row
}
