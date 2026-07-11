package builders

import (
	"fmt"
	"io"
	"net/url"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HTMXInlineFormBuilder struct {
	endpoint     url.URL
	submitMethod HTMXMethod
	numColumns   int
	fields       []FormField
}

func NewHTMXInlineFormBuilder(numColumns int, endpoint url.URL, submitMethod HTMXMethod) HTMXInlineFormBuilder {
	return HTMXInlineFormBuilder{
		endpoint:     endpoint,
		numColumns:   numColumns,
		submitMethod: submitMethod,
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
				Class("flex p-2"),
				Attr(fmt.Sprintf(`%v=%v`, b.submitMethod.String(), b.endpoint.String())),
				Attr(`hx-target="closest tr"`),
				Attr(`hx-swap=outerHTML`),
				Div(
					Class("ml-auto flex items-center gap-2"),
					Map(b.fields, func(f FormField) Node {
						return Div(
							Class("flex flex-col gap-1"),
							Group([]Node{f.Label(), f.Input()}))
					}),
					Div(Class("ml-auto flex gap-2"),
						Button(Type("submit"),
							Text("Submit"),
							Class(st.SuccessButton)),
						cmp.MakeHTMXAbortButton("Cancel", b.endpoint),
					),
				),
			)),
	)

	form.Render(w)

}
