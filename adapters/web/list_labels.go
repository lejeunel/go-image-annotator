package web

import (
	"io"
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	. "maragu.dev/gomponents"
)

type ListLabelsPresenter struct {
	b.PageBuilder
	Writer io.Writer
	WebPageErrorPresenter
}

func (p ListLabelsPresenter) Success(r list.Response) {
	listBuilder := b.NewPaginatedListBuilder([]string{"name", "description", "actions"}, rt.Labels, r.Pagination)
	for _, l := range r.Labels {
		actions := b.NewActionsPanelBuilder()
		editURL, _ := url.Parse("/edit-url")
		deleteURL, _ := url.Parse("/delete-url")
		actions.SetEdit(*editURL)
		actions.SetConfirmDelete(*deleteURL)
		row := tb.NewRow()
		row.AddCell(tb.NewCell(Text(l.Name)))
		row.AddCell(tb.NewCell(Text(l.Description)))
		row.AddCell(tb.NewCell(actions.Build()))
		listBuilder.AddRow(row)
	}
	p.PageBuilder.SetContent(listBuilder.Build(), nil)
	p.Render(p.Writer)
}

func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.Label.List.Execute(r.Context(),
		list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListLabelsPresenter(w, s.PageBuilder))
}

func NewListLabelsPresenter(w http.ResponseWriter, p b.PageBuilder) ListLabelsPresenter {
	return ListLabelsPresenter{
		*p.SetTitle("Labels").SetActive(cmp.LabelsPageActive),
		w, NewWebPageErrorPresenter(w)}
}
