package web

import (
	"net/http"
	"net/url"

	html "github.com/lejeunel/go-image-annotator-v2/shared/html"
	n "github.com/lejeunel/go-image-annotator-v2/shared/navigation"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/label/list"
	. "maragu.dev/gomponents"
)

type ListLabelsPresenter struct {
	ListRenderer
}

func (p ListLabelsPresenter) Success(r list.Response) {
	table := html.PaginationTable{Fields: []string{"name", "description"}}
	for _, l := range r.Labels {
		table.Rows = append(table.Rows,
			html.PaginationTableRow{Values: []Node{Text(l.Name), Raw(l.Description)}})
	}
	p.RenderSuccess(table, r.Pagination)
}

func (s *Server) ListLabels(w http.ResponseWriter, r *http.Request) {
	s.Label.List.Execute(list.Request{PageSize: s.Label.DefaultPageSize, Page: int64(GetPageFromRequest(r))},
		NewListLabelsPresenter(w))
}

func NewListLabelsPresenter(w http.ResponseWriter) ListLabelsPresenter {
	baseURL, _ := url.Parse("/labels")
	return ListLabelsPresenter{
		ListRenderer: NewListRenderer("Labels", *baseURL,
			n.LabelsPageActive, w),
	}
}
