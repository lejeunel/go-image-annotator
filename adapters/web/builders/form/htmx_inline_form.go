package form

import (
	"bytes"
	"fmt"
	"io"
	"net/url"

	cmp "github.com/lejeunel/go-image-annotator/adapters/web/components"
	st "github.com/lejeunel/go-image-annotator/adapters/web/styles"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type FormMode int

const (
	EditMode FormMode = iota
	CreateMode
)

func (m FormMode) Verb() string {
	switch m {
	case EditMode:
		return "Editing"
	case CreateMode:
		return "Creating"
	default:
		return ""
	}
}
func (m FormMode) HTMXMethod() string {
	switch m {
	case EditMode:
		return "hx-put"
	case CreateMode:
		return "hx-post"
	default:
		return ""
	}
}

type HTMXInlineFormBuilder struct {
	resourceName string
	endpoint     url.URL
	mode         FormMode
	numColumns   int
	title        *string
	fields       []Renderer
}
type FormOption func(*HTMXInlineFormBuilder)

func WithMode(m FormMode) FormOption {
	return func(b *HTMXInlineFormBuilder) {
		b.mode = m
	}
}

func NewHTMXInlineFormBuilder(resourceName string, numColumns int, endpoint url.URL, opts ...FormOption) HTMXInlineFormBuilder {
	f := &HTMXInlineFormBuilder{
		resourceName: resourceName,
		endpoint:     endpoint,
		numColumns:   numColumns,
		mode:         EditMode,
	}

	for _, opt := range opts {
		opt(f)
	}
	return *f
}
func (b *HTMXInlineFormBuilder) AddTitle(title string) *HTMXInlineFormBuilder {
	b.title = &title
	return b
}
func (b *HTMXInlineFormBuilder) AddTextField(fieldName, displayName, divId string, opts ...FormTextFieldOption) *HTMXInlineFormBuilder {
	field := NewFormTextField(fieldName, displayName, divId, opts...)
	b.fields = append(b.fields, field)
	return b
}
func (b *HTMXInlineFormBuilder) AddSelectableCombobox(title, id string) *SelectableCombobox {
	box := NewSelectableCombobox(title, id)
	b.fields = append(b.fields, &box)
	return &box
}
func (b HTMXInlineFormBuilder) Render(w io.Writer) {

	caption := Div(
		Class("ml-auto flex gap-2"),
		Div(Text(b.mode.Verb())),
		Div(Class("font-bold"), Text(b.resourceName)))
	form := Tr(
		Td(Attr(fmt.Sprintf("colspan=%v", b.numColumns)),
			Form(
				Class("flex p-2"),
				Attr(fmt.Sprintf(`%v=%v`, b.mode.HTMXMethod(), b.endpoint.String())),
				Attr(`hx-target="closest tr"`),
				Attr(`hx-swap=outerHTML`),
				Div(
					caption,
					Class("ml-auto flex items-center gap-2"),
					Map(b.fields, func(f Renderer) Node {
						var buf bytes.Buffer
						f.Render(&buf)
						return Div(
							Class("flex flex-col gap-1"),
							Raw(buf.String()))
					}),
					Div(Class("ml-auto flex gap-2"),
						Button(Type("submit"),
							Text("Submit"),
							Class(st.SuccessButton)),
						cmp.MakeHTMXAbortButton("Cancel", b.endpoint.String()),
					),
				),
			)),
	)

	form.Render(w)

}
