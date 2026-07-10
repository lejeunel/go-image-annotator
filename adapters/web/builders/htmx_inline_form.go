package builders

import (
	"fmt"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HTMXInlineFormBuilder struct {
	submitEndpoint string
	submitMethod   HTMXMethod
	numColumns     int
	fields         []FormField
}

func NewHTMXInlineFormBuilder(numColumns int, submitEndpoint string, submitMethod HTMXMethod) HTMXInlineFormBuilder {
	return HTMXInlineFormBuilder{
		submitEndpoint: submitEndpoint,
		numColumns:     numColumns,
		submitMethod:   submitMethod,
	}
}

func (b *HTMXInlineFormBuilder) AddTextField(fieldName, displayName, divId string, opts ...FormTextFieldOption) *HTMXInlineFormBuilder {
	field := NewFormTextField(fieldName, displayName, divId, opts...)
	b.fields = append(b.fields, field)
	return b
}

func (b HTMXInlineFormBuilder) Render(w io.Writer) {
	form := Tr(
		Td(Attr(fmt.Sprintf("colspan=%v", b.numColumns)),
			Form(
				Class("flex items-end gap-4 p-2"),
				Attr(fmt.Sprintf(`%v=%v`, b.submitMethod.String(), b.submitEndpoint)),
				Attr(`hx-target="closest tr"`),
				Attr(`hx-swap=outerHTML`),
				Class("flex flex-col gap-1"),
				Map(b.fields, func(f FormField) Node {
					return Div(
						Class("flex flex-col gap-1"),
						Group([]Node{f.Label(), f.Input()}))
				}),
				Div(Class("ml-auto flex gap-2"),
					Button(Type("submit"),
						Text("Submit"),
						Class(st.SuccessButton)),
					Button(Type("button"),
						Text("Cancel"),
						Class(st.InactiveButton)),
				),
			)),
	)

	form.Render(w)

}
