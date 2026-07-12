package web

import (
	"io"
	"net/http"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
)

var listCollectionsFields = []string{"name", "description", "group", "created", "actions"}

type ListCollectionsPresenter struct {
	b.PaginatedListBuilder
	Writer io.Writer
	WebPageErrorPresenter
}

func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	p.SetTitle("Collections").SetActive(cmp.CollectionsPageActive)
	b := b.NewPaginatedListBuilder(p, listCollectionsFields)
	return ListCollectionsPresenter{b, w, NewWebPageErrorPresenter(w)}
}
func (p ListCollectionsPresenter) SuccessFindCollection(c clc.Collection) {
	MakeListCollectionRow(c).Render(p.Writer)
}
func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	p.SetPagination(r.Pagination, rt.Collections)
	for _, c := range r.Collections {
		row := MakeListCollectionRow(c)
		p.AddRow(row)
	}
	p.AddCreationButton("Create new collection", rt.CreateCollectionForm, createCollectionTargetDiv)
	p.Build().Render(p.Writer)
}

func (s *Server) GetCollection(w http.ResponseWriter, r *http.Request) {
	s.Collection.Find.Execute(r.Context(),
		r.URL.Query().Get("name"),
		NewListCollectionsPresenter(w, s.PageBuilder))
}
func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentity(r.Context())
	s.Collection.List.Execute(r.Context(), list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w, s.PageBuilder))
}

func MakeListCollectionRow(c clc.Collection) tb.Row {
	var groupName string
	if c.Group == nil {
		groupName = "n/a"
	} else {
		groupName = c.Group.Name
	}

	actions := b.NewActionsPanelBuilder()
	actions.SetEdit(rt.AddQueryParams(rt.EditCollectionForm,
		"name", c.Name, "description", c.Description))
	actions.SetConfirmDelete(rt.AddQueryParams(rt.ConfirmDeleteCollection, "name", c.Name))

	row := tb.NewRow()
	row.AddCell(tb.NewCell(cmp.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name)))
	row.AddCell(tb.NewCell(Text(c.Description)))
	row.AddCell(tb.NewCell(Text(groupName)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(c.CreatedAt))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
