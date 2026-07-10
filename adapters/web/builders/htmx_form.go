package builders

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

func (b *FormBuilder) AddTextField(fieldName, displayName, divId string, required bool) *FormBuilder {
	field := NewFormTextField(fieldName, displayName, divId, WithRequired())
	b.fields = append(b.fields, field)
	return b
}

func NewFormBuilder(submitEndpoint string) FormBuilder {
	return FormBuilder{submitEndpoint: submitEndpoint}
}
