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
	fields         []FormField
}

func (b *FormBuilder) AddTextField(fieldName, displayName, divId string, opts ...FormTextFieldOption) *FormBuilder {
	field := NewFormTextField(fieldName, displayName, divId, opts...)
	b.fields = append(b.fields, field)
	return b
}
func (b *FormBuilder) AddCheckbox(fieldName, displayName, divId string) *FormBuilder {
	field := NewFormCheckboxField(fieldName, displayName, divId)
	b.fields = append(b.fields, field)
	return b
}

func NewFormBuilder(submitEndpoint string) FormBuilder {
	return FormBuilder{submitEndpoint: submitEndpoint}
}
