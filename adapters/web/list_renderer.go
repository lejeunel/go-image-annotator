package web

import (
	"io"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	html "github.com/lejeunel/go-image-annotator/adapters/web/html"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/url"
)

type ListRenderer struct {
	b.PageBuilder
	ListURL url.URL
	Writer  io.Writer
}

func (p ListRenderer) RenderSuccess(table html.MyTable, pagination pagination.Pagination, header *Node) {

	content := html.MakePaginatedContent(p.ListURL, table, pagination)
	var cat Node
	if header != nil {
		cat = Div(*header, content)
	} else {
		cat = content
	}
	p.PageBuilder.SetContent(cat).Render(p.Writer)
}

func (p ListRenderer) Error(err error) {
	p.PageBuilder.SetError(err).Render(p.Writer)
}

func NewListRenderer(pageBuilder b.PageBuilder, listURL url.URL, w io.Writer) ListRenderer {
	return ListRenderer{PageBuilder: pageBuilder, ListURL: listURL,
		Writer: w}

}
