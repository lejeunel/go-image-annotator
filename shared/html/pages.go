package html

import (
	"net/url"

	n "github.com/lejeunel/go-image-annotator/shared/navigation"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakePaginatedContent(baseURL url.URL, table PaginationTable, p pagination.Pagination) Node {
	paginator := MakePaginator(baseURL, int(p.Page), int(p.TotalPages), len(table.Rows), int(p.TotalRecords))
	return Div(Div(Class("py-2"), paginator), table.Build())

}

func MakePaginatedView(baseURL url.URL, title string, pagination pagination.Pagination,
	table PaginationTable, activePage n.ActivePage) Node {

	content := MakePaginatedContent(baseURL, table, pagination)
	p := NewTitledPageBuilder(title)
	p.SetContent(content)
	p.SetActive(activePage)
	return p.Build()
}
