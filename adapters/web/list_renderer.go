package web

import (
	"io"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	"net/url"
)

type ListRenderer struct {
	b.PageBuilder
	ActivePage b.ActivePage
	ListURL    url.URL
	Writer     io.Writer
}

func (p ListRenderer) RenderSuccess(table html.PaginationTable, pagination pagination.Pagination) {

	content := html.MakePaginatedContent(p.ListURL, table, pagination)
	p.PageBuilder.SetContent(content).SetActive(p.ActivePage).Render(p.Writer)
}

func (p ListRenderer) Error(err error) {
	b.NewPageBuilder(p.APIPath).SetError(err).Render(p.Writer)
}

func NewListRenderer(pageBuilder b.PageBuilder, listURL url.URL, page b.ActivePage, w io.Writer) ListRenderer {
	return ListRenderer{PageBuilder: pageBuilder, ListURL: listURL,
		ActivePage: page, Writer: w}

}
