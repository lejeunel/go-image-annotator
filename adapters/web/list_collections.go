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
	table := html.PaginationTable{Fields: []string{"name", "description", "created"}}
	for _, c := range r.Collections {
		table.Rows = append(table.Rows,
			html.PaginationTableRow{Values: []Node{html.MakeTextLink("/images?collection="+c.Name, c.Name),
				Raw(c.Description), Raw(DateTimeToStr(c.CreatedAt))}})
	}
	p.RenderSuccess(table, r.Pagination)
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.Collection.List.Execute(list.Request{PageSize: s.Collection.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w))
}

func NewListCollectionsPresenter(w http.ResponseWriter) ListCollectionsPresenter {
	baseURL, _ := url.Parse("/collections")
	return ListCollectionsPresenter{
		ListRenderer: NewListRenderer("Collections", *baseURL,
			n.CollectionsPageActive, w),
	}
}
