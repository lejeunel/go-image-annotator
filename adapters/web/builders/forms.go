package builders

import (
	"fmt"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type CreateFormTextField struct {
	fieldName   string
	displayName string
	divId       string
	required    bool
}

func (f CreateFormTextField) Label() Node {
	displayName := f.displayName
	if !f.required {
		displayName = displayName + " (Optional)"
	}
	return Label(For(f.fieldName), Text(displayName), Class("block text-sm font-medium text-gray-900 dark:text-white"))
}

func (f CreateFormTextField) Input() Node {
	return Input(Type("text"), ID(f.divId), Name(f.fieldName), If(f.required, Required()),
		Class("w-full px-3 py-2 mb-6 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"))
}

func NewCreateFormTextField(fieldName, displayName, divId string, required bool) CreateFormTextField {
	return CreateFormTextField{fieldName, displayName, divId, required}
}

type FormField interface {
	Label() Node
	Input() Node
}
type CreateFormBuilder struct {
	postRoute   string
	containerId string
	fields      []FormField
}

func NewCreateFormBuilder(postRoute string, containerId string) CreateFormBuilder {
	return CreateFormBuilder{postRoute: postRoute, containerId: containerId}
}
func (b *CreateFormBuilder) AddTextField(fieldName, displayName, divId string, required bool) *CreateFormBuilder {
	field := NewCreateFormTextField(fieldName, displayName, divId, required)
	b.fields = append(b.fields, field)
	return b
}

func (b CreateFormBuilder) Render(w io.Writer) {
	form := Span(Class("w-full inline-flex items-center justify-start mt-2"),
		Form(
			Attr(fmt.Sprintf(`hx-post=%v`, b.postRoute)),
			Attr(`hx-swap="none"`),
			Class("bg-white dark:bg-gray-800 p-8 rounded-lg shadow-md w-80 mb-4"),
			Map(b.fields, func(f FormField) Node {
				return Group([]Node{f.Label(), f.Input()})
			}),
			Span(Class("flex items-center gap-2"),
				Button(Type("submit"),
					Text("Submit"),
					Class(st.PrimaryButton)),
				Button(Type("button"),
					Text("Cancel"),
					Class(st.InactiveButton),
					Attr(`hx-on:click`, fmt.Sprintf(`document.getElementById('%v').innerHTML=''`, b.containerId))),
			),
		))

	form.Render(w)

}
