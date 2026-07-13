package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	tb "github.com/lejeunel/go-image-annotator/adapters/web/builders/table"
	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HTMXErrorPresenter struct {
	task   string
	writer http.ResponseWriter
}

func (p HTMXErrorPresenter) Error(err error) {
	payload, _ := NotifyError(fmt.Sprintf("failed %v", p.task),
		err.Error())
	fmt.Println(err)
	fmt.Println(payload)
	p.writer.Header().Set("HX-Trigger", string(payload))
	p.writer.WriteHeader(http.StatusUnprocessableEntity)
}

func NewHTMXErrorPresenter(t string, w http.ResponseWriter) HTMXErrorPresenter {
	return HTMXErrorPresenter{task: t, writer: w}
}

type WebPageErrorPresenter struct {
	b.PageBuilder
	writer http.ResponseWriter
}

func (p WebPageErrorPresenter) Error(err error) {
	p.PageBuilder.SetError(err).Render(p.writer)
}

func NewWebPageErrorPresenter(w http.ResponseWriter) WebPageErrorPresenter {
	return WebPageErrorPresenter{writer: w}
}

func RenderConfirmDeleteRow(numCols int, name, resourceType string,
	deleteURL url.URL, cancelURL url.URL, w io.Writer) {
	row := tb.NewRow()
	row.AddCell(
		tb.NewCell(
			Div(Text("Do you really want to delete "+resourceType),
				Div(Class("font-bold"), Text(fmt.Sprintf("%v?", name)))),
			tb.WithCellAttr(fmt.Sprintf(`colspan=%v`, numCols-1)),
			tb.WithCellClass("text-right"),
		))
	row.AddCell(tb.NewCell(Span(Class("flex items-center gap-2"),
		cmp.MakeHTMXDeleteButton("Yes", deleteURL.String()),
		cmp.MakeHTMXAbortButton("Cancel", cancelURL.String()))))
	row.Render(w)
}
