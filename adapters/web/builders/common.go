package builders

import (
	"fmt"
	"io"
	"net/url"

	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func RenderConfirmDeleteRow(numCols int, name, resourceType string,
	url url.URL, w io.Writer) {
	row := tb.NewRow()
	row.AddCell(
		tb.NewCell(
			Div(Text("Do you really want to delete "+resourceType),
				Div(Class("font-bold"), Text(fmt.Sprintf("%v?", name)))),
			tb.WithCellAttr(fmt.Sprintf(`colspan=%v`, numCols-1)),
			tb.WithCellClass("text-right"),
		))
	row.AddCell(tb.NewCell(Span(Class("flex items-center gap-2"),
		cmp.MakeHTMXDeleteButton("Yes", url.String()),
		cmp.MakeHTMXAbortButton("Cancel", url.String()))))
	row.Render(w)
}
