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

type HTMXFormBuilder struct {
	containerId string
	buttonAttrs []string
	title       *string
	FormBuilder
}

func NewHTMXFormBuilder(submitEndpoint string, containerId string) HTMXFormBuilder {
	return HTMXFormBuilder{
		FormBuilder: FormBuilder{submitEndpoint: submitEndpoint},
		containerId: containerId,
	}
}

func (b *HTMXFormBuilder) AddSubmitQueryParam(key, value string) *HTMXFormBuilder {
	url := rt.AddQueryParams(b.FormBuilder.submitEndpoint, key, value)
	b.FormBuilder.submitEndpoint = url.String()
	return b
}

func (b *HTMXFormBuilder) AddTitle(title string) *HTMXFormBuilder {
	b.title = &title
	return b
}

func (b *HTMXFormBuilder) AddButtonAttr(attr string) *HTMXFormBuilder {
	b.buttonAttrs = append(b.buttonAttrs, attr)
	return b
}

func (b HTMXFormBuilder) Build() Node {
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
				"w-120 bg-surface-alt dark:bg-surface-dark-alt border-outline dark:border-outline-dark p-4 rounded-lg shadow-md mb-4",
			),
			title,
			Map(b.FormBuilder.fields, func(f Renderer) Node {
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

func (b HTMXFormBuilder) Render(w io.Writer) {
	b.Build().Render(w)
}
