package web

import (
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	. "maragu.dev/gomponents"
)

type ListLabelsPresenter struct {
	ListRenderer
}

func (p ListLabelsPresenter) Success(r list.Response) {
	table := html.MyTable{Fields: []string{"name", "description"}}
	for _, l := range r.Labels {
		table.Rows = append(table.Rows,
			html.MyTableRow{Values: []Node{Text(l.Name), Raw(l.Description)}})
	}
	p.RenderSuccess(table, r.Pagination, nil)
}

func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request) {
	s.PageBuilder.SetUserIdentityFromContext(r.Context())
	s.Label.List.Execute(r.Context(),
		list.Request{PageSize: s.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListLabelsPresenter(w, s.PageBuilder))
}

func NewListLabelsPresenter(w http.ResponseWriter, p b.PageBuilder) ListLabelsPresenter {
	baseURL, _ := url.Parse("/labels")
	return ListLabelsPresenter{
		ListRenderer: NewListRenderer(*p.SetTitle("Labels").SetActive(b.LabelsPageActive), *baseURL,
			w),
	}
}
