package components

import (
	"fmt"
	"net/url"

	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	s "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakeHTMXCreateButton(text string, hxPut string, hxTarget string) Node {
	return Div(
		Class("m-2"),
		Button(
			Attr(fmt.Sprintf("hx-get=%v", hxPut)),
			Attr(fmt.Sprintf("hx-target=#%v", hxTarget)),
			Attr(`hx-swap=innerHTML`),
			Class(s.PrimaryButton),
			Span(Class("flex items-center"),
				Raw(ic.AddIcon), Div(Class("p-1"), Text(text)),
			)))
}

func MakeHTMXDeleteButton(text string, hxDelete url.URL) Node {
	return Div(
		Button(
			Attr(fmt.Sprintf("hx-delete=%v", hxDelete.String())),
			Attr(`hx-include="closest tr"`),
			Class(s.DangerButton),
			Div(Class("p-1"), Text(text))))
}

func MakeHTMXAbortButton(text string, hxGet url.URL) Node {
	return Div(
		Button(
			Attr(fmt.Sprintf("hx-get=%v", hxGet.String())),
			Class(s.InvertButton),
			Div(Class("p-1"), Text(text))))
}
