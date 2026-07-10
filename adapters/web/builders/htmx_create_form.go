package builders

import (
	"fmt"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HTMXCreateFormBuilder struct {
	containerId string
	FormBuilder
}

func NewHTMXCreateFormBuilder(submitEndpoint string, containerId string) HTMXCreateFormBuilder {
	return HTMXCreateFormBuilder{
		FormBuilder: FormBuilder{submitEndpoint: submitEndpoint},
		containerId: containerId}
}
func (b *HTMXCreateFormBuilder) AddTextField(fieldName, displayName, divId string, required bool) *HTMXCreateFormBuilder {
	field := NewFormTextField(fieldName, displayName, divId, WithRequired())
	b.fields = append(b.fields, field)
	return b
}

func (b HTMXCreateFormBuilder) Render(w io.Writer) {
	form := Span(Class("w-full inline-flex items-center justify-start mt-2"),
		Form(
			Attr(fmt.Sprintf(`hx-post=%v`, b.submitEndpoint)),
			Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
			Map(b.fields, func(f FormField) Node {
				return Group([]Node{f.Label(), f.Input()})
			}),
			Span(Class("flex items-center gap-2"),
				Button(Type("submit"),
					Text("Submit"),
					Class(st.SuccessButton)),
				Button(Type("button"),
					Text("Cancel"),
					Class(st.InactiveButton),
					Attr(`hx-on:click`, fmt.Sprintf(`document.getElementById('%v').innerHTML=''`, b.containerId))),
			),
		))

	form.Render(w)

}
