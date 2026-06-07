package web

import (
	"net/http"
	"net/url"

	html "github.com/lejeunel/go-image-annotator/shared/html"
	n "github.com/lejeunel/go-image-annotator/shared/navigation"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
)

type ListCollectionsPresenter struct {
	ListRenderer
}

func (p ListCollectionsPresenter) Success(r list.Response) {
	table := html.PaginationTable{Fields: []string{"name", "description", "group", "created"}}
	for _, c := range r.Collections {
		table.Rows = append(table.Rows,
			html.PaginationTableRow{Values: []Node{html.MakeTextLink("/images?collection="+c.Name, c.Name),
				Raw(c.Description), Raw(c.Group), Raw(DateTimeToStr(c.CreatedAt))}})
	}
	p.RenderSuccess(table, r.Pagination)
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	s.Collection.List.Execute(r.Context(), list.Request{PageSize: s.Collection.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w, s.PageBuilder))
}

func NewListCollectionsPresenter(w http.ResponseWriter, p html.PageBuilder) ListCollectionsPresenter {
	baseURL, _ := url.Parse("/collections")
	return ListCollectionsPresenter{
		ListRenderer: NewListRenderer(*p.SetTitle("Collections"), *baseURL,
			n.CollectionsPageActive, w),
	}
}
