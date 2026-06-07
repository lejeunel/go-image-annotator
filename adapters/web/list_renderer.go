package web

import (
	"io"

	html "github.com/lejeunel/go-image-annotator/shared/html"
	n "github.com/lejeunel/go-image-annotator/shared/navigation"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	"net/url"
)

type ListRenderer struct {
	html.PageBuilder
	ActivePage n.ActivePage
	ListURL    url.URL
	Writer     io.Writer
}

func (p ListRenderer) RenderSuccess(table html.PaginationTable, pagination pagination.Pagination) {

	content := html.MakePaginatedContent(p.ListURL, table, pagination)
	p.PageBuilder.SetContent(content).SetActive(p.ActivePage).Render(p.Writer)
}

func (p ListRenderer) Error(err error) {
	html.NewPageBuilder(p.APIPath).SetError(err).Render(p.Writer)
}

func NewListRenderer(pageBuilder html.PageBuilder, listURL url.URL, page n.ActivePage, w io.Writer) ListRenderer {
	return ListRenderer{PageBuilder: pageBuilder, ListURL: listURL,
		ActivePage: page, Writer: w}

}
