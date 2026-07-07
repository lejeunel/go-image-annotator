package web

import (
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	. "maragu.dev/gomponents"
)

type ListCollectionsPresenter struct {
	ListRenderer
}

func (p ListCollectionsPresenter) Success(r list.Response) {
	table := html.MyTable{Fields: []string{"name", "description", "group", "created"}}
	for _, c := range r.Collections {
		var groupName string
		if c.Group == nil {
			groupName = "n/a"
		} else {
			groupName = c.Group.Name
		}
		table.Rows = append(table.Rows,
			html.MyTableRow{Values: []Node{html.MakeTextLink("/images?collection="+c.Name, c.Name),
				Raw(c.Description), Raw(groupName), Raw(DateTimeToStr(c.CreatedAt))}})
	}

	p.RenderSuccess(table, r.Pagination, nil)
}

func (s *Server) ListCollections(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	s.Collection.List.Execute(r.Context(), list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListCollectionsPresenter(w, s.PageBuilder))
}

func NewListCollectionsPresenter(w http.ResponseWriter, p b.PageBuilder) ListCollectionsPresenter {
	baseURL, _ := url.Parse("/collections")
	return ListCollectionsPresenter{
		ListRenderer: NewListRenderer(*p.SetTitle("Collections").SetActive(b.CollectionsPageActive), *baseURL,
			w),
	}
}
