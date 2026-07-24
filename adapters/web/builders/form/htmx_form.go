package form

type HTMXMethod int

const (
	HTMXPostMethod HTMXMethod = iota
	HTMXPutMethod
)

func (m HTMXMethod) String() string {
	switch m {
	case HTMXPostMethod:
		return "hx-post"
	case HTMXPutMethod:
		return "hx-put"
	default:
		return "hx-post"
	}
}

type FormBuilder struct {
	submitEndpoint string
	fields         []Renderer
}

func (b *FormBuilder) AddTextField(fieldName, displayName string, opts ...FormTextFieldOption) *FormBuilder {
	field := NewFormTextField(fieldName, displayName, opts...)
	b.fields = append(b.fields, field)
	return b
}
func (b *FormBuilder) AddCheckbox(fieldName, displayName string) *FormBuilder {
	field := NewFormCheckboxField(fieldName, displayName)
	b.fields = append(b.fields, field)
	return b
}

func NewFormBuilder(submitEndpoint string) FormBuilder {
	return FormBuilder{submitEndpoint: submitEndpoint}
}
