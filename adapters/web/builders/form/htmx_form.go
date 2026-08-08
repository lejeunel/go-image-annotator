package form

type FormBuilder struct {
	submitEndpoint string
	fields         []Renderer
}

func (b *FormBuilder) AddTextField(
	fieldName, displayName string,
	opts ...FormFieldOption,
) *FormBuilder {
	field := NewFormTextField(fieldName, displayName, opts...)
	b.fields = append(b.fields, field)
	return b
}

func (b *FormBuilder) AddIntField(
	fieldName, displayName string,
	opts ...FormFieldOption,
) *FormBuilder {
	field := NewFormIntField(fieldName, displayName, opts...)
	b.fields = append(b.fields, field)
	return b
}

func (b *FormBuilder) AddFloatField(
	fieldName, displayName string,
	opts ...FormFieldOption,
) *FormBuilder {
	field := NewFormFloatField(fieldName, displayName, opts...)
	b.fields = append(b.fields, field)
	return b
}

func (b *FormBuilder) AddDateTimeField(
	fieldName, displayName string,
	opts ...FormFieldOption,
) *FormBuilder {
	field := NewFormDateTimeField(fieldName, displayName, opts...)
	b.fields = append(b.fields, field)
	return b
}

func (b *FormBuilder) AddBooleanField(
	fieldName, displayName string,
	opts ...FormFieldOption,
) *FormBuilder {
	field := NewFormBooleanField(fieldName, displayName, opts...)
	b.fields = append(b.fields, field)
	return b
}

func (b *FormBuilder) AddCheckbox(fieldName, displayName string) *FormBuilder {
	field := NewFormCheckboxField(fieldName, displayName)
	b.fields = append(b.fields, field)
	return b
}
