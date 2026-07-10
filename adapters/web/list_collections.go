package web

import (
	"io"
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
)

type ListCollectionsPresenter struct {
	b.PageBuilder
	Writer io.Writer
	WebPageErrorPresenter
}

var listCollectionsFields = []string{"name", "description", "group", "created", "actions"}

func MakeListCollectionRow(c clc.Collection) tb.Row {
	var groupName string
	if c.Group == nil {
		groupName = "n/a"
	} else {
		groupName = c.Group.Name
	}

	actions := b.NewActionsPanelBuilder()

	editURL, _ := url.Parse("/edit-url")
	actions.SetEdit(*editURL)
	actions.SetConfirmDelete(rt.AppendValueToQueryArgs(rt.ConfirmDeleteCollection, "name", c.Name))
	row := tb.NewRow()
	row.AddCell(tb.NewCell(html.MakeTextLink(rt.MakeImagesURL(c.Name), c.Name)))
	row.AddCell(tb.NewCell(Text(c.Description)))
	row.AddCell(tb.NewCell(Text(groupName)))
	row.AddCell(tb.NewCell(Text(cmp.DateTimeToStr(c.CreatedAt))))
	row.AddCell(tb.NewCell(actions.Build()))
	return row

}
func (p ListCollectionsPresenter) SuccessFindCollection(c clc.Collection) {
	MakeListCollectionRow(c).Render(p.Writer)
}
func (p ListCollectionsPresenter) SuccessListCollections(r list.Response) {
	listBuilder := b.NewPaginatedListBuilder(
		listCollectionsFields,
		rt.Collections,
		r.Pagination)
	for _, c := range r.Collections {
		row := MakeListCollectionRow(c)
		listBuilder.AddRow(row)
	}
	listBuilder.AddCreationButton("Create new collection", rt.CreateCollectionForm, createCollectionTargetDiv)
	p.PageBuilder.SetContent(listBuilder.Build(), nil)
	p.Render(p.Writer)
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
func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	return ListCollectionsPresenter{*p.SetTitle("Collections").SetActive(cmp.CollectionsPageActive), w,
		NewWebPageErrorPresenter(w)}
}
