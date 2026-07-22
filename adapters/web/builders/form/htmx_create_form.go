package form

import (
	"fmt"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HTMXCreateFormBuilder struct {
	containerId string
	title       *string
	FormBuilder
}

func NewHTMXCreateFormBuilder(submitEndpoint string, containerId string) HTMXCreateFormBuilder {
	return HTMXCreateFormBuilder{
		FormBuilder: FormBuilder{submitEndpoint: submitEndpoint},
		containerId: containerId}
}

func (b *HTMXCreateFormBuilder) AddTitle(title string) *HTMXCreateFormBuilder {
	b.title = &title
	return b
}
func (b HTMXCreateFormBuilder) Render(w io.Writer) {
	var title Node
	if b.title != nil {
		title = Div(Class("ml-auto flex gap-2 font-bold text-lg"),
			Text(*b.title))
	}

	form := Span(Class("w-full inline-flex items-center justify-start mt-2"),
		Form(
			Attr(fmt.Sprintf(`hx-post=%v`, b.submitEndpoint)),
			Class("bg-surface-alt/50 dark:bg-surface-dark-alt/50 p-8 rounded-lg shadow-md w-80 mb-4"),
			title,
			Map(b.fields, func(f FormField) Node {
				return Group([]Node{Div(Class("mb-3"), f.Build())})
			}),
			Span(Class("flex items-center gap-2"),
				Button(Type("submit"),
					Text("Submit"),
					Class(st.SuccessButton)),
				Button(Type("button"),
					Text("Cancel"),
					Class(st.AbortButton),
					Attr(`hx-on:click`, fmt.Sprintf(`document.getElementById('%v').innerHTML=''`, b.containerId))),
			),
		))

	form.Render(w)

}
