package form

import (
	"bytes"
	"fmt"
	"io"

	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	rt "github.com/lejeunel/go-image-annotator/routes"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type HTMXCreateFormBuilder struct {
	containerId string
	buttonAttrs []string
	title       *string
	FormBuilder
}

func NewHTMXCreateFormBuilder(submitEndpoint string, containerId string) HTMXCreateFormBuilder {
	return HTMXCreateFormBuilder{
		FormBuilder: FormBuilder{submitEndpoint: submitEndpoint},
		containerId: containerId,
	}
}
func (b *HTMXCreateFormBuilder) AddSubmitQueryParam(key, value string) *HTMXCreateFormBuilder {
	url := rt.AddQueryParams(b.FormBuilder.submitEndpoint, key, value)
	b.FormBuilder.submitEndpoint = url.String()
	return b
}

func (b *HTMXCreateFormBuilder) AddTitle(title string) *HTMXCreateFormBuilder {
	b.title = &title
	return b
}

func (b *HTMXCreateFormBuilder) AddButtonAttr(attr string) *HTMXCreateFormBuilder {
	b.buttonAttrs = append(b.buttonAttrs, attr)
	return b
}

func (b HTMXCreateFormBuilder) Build() Node {
	var title Node
	if b.title != nil {
		title = Div(Class("ml-auto flex gap-2 font-bold"),
			Text(*b.title))
	}

	attrs := []Node{Attr(fmt.Sprintf(`hx-post=%v`, b.submitEndpoint))}
	for _, a := range b.buttonAttrs {
		attrs = append(attrs, Attr(a))
	}
	return Span(Class("w-full inline-flex items-center justify-start mt-1"),
		Form(
			Group(attrs),
			Class(
				"bg-surface-alt/50 dark:bg-surface-dark-alt/50 p-4 rounded-lg shadow-md w-80 mb-4",
			),
			title,
			Map(b.fields, func(f Renderer) Node {
				var buf bytes.Buffer
				f.Render(&buf)
				return Group([]Node{Div(Class("mb-3"), Raw(buf.String()))})
			}),
			Span(Class("flex items-center gap-2"),
				Button(Type("submit"),
					Text("Submit"),
					Class(st.SuccessButton)),
				Button(
					Type("button"),
					Text("Cancel"),
					Class(st.AbortButton),
					Attr(
						`hx-on:click`,
						fmt.Sprintf(`document.getElementById('%v').innerHTML=''`, b.containerId),
					),
				),
			),
		))
}

func (b HTMXCreateFormBuilder) Render(w io.Writer) {
	b.Build().Render(w)
}
