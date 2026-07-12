package components

import (
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakePaginatedContent(baseURL string, numRows int, table Node, p pagination.Pagination) Node {
	paginator := MakePaginator(baseURL, int(p.Page), int(p.TotalPages), numRows, int(p.TotalRecords))
	return Div(Div(Class("py-2"), paginator), table)

}
