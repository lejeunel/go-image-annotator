package html

import (
	"net/url"

	"github.com/lejeunel/go-image-annotator/shared/pagination"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakePaginatedContent(baseURL url.URL, table MyTable, p pagination.Pagination) Node {
	paginator := MakePaginator(baseURL, int(p.Page), int(p.TotalPages), len(table.Rows), int(p.TotalRecords))
	return Div(Div(Class("py-2"), paginator), table.Build())

}
