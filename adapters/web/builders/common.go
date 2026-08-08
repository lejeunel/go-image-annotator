package builders

import (
	"fmt"
	"io"
	"net/url"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func RenderConfirmDeleteRow(numCols int, name, resourceType string,
	url url.URL, w io.Writer,
) {
	content := Tr(
		Td(
			Attr(fmt.Sprintf("colspan=\"%v\"", numCols)),
			Div(Class("flex items-center justify-end py-2 mr-2"),
				Div(Class("mr-3 "), Text("Do you really want to delete "+resourceType),
					Span(Class("ml-2 font-bold"), Text(fmt.Sprintf("%v?", name)))),
				Span(Class("flex items-center gap-2"),
					cmp.MakeHTMXDeleteButton("Yes", url.String()),
					cmp.MakeHTMXAbortButton("Cancel", url.String())),
			)))
	content.Render(w)
}
