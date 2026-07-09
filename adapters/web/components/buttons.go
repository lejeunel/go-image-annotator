package components

import (
	"fmt"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	s "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MakeHTMXCreateButton(text string, hxGet string, hxTarget string) Node {
	return Button(
		Attr(fmt.Sprintf("hx-get=%v", hxGet)),
		Attr(fmt.Sprintf("hx-target=#%v", hxTarget)),
		Attr(`hx-swap=innerHTML`),
		Class(s.ActivePrimaryButton),
		Raw(ic.AddIcon), Div(Class("p-1"), Text("Create new collection")))

}
