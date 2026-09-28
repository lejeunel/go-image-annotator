package form

import (
	"io"
	"strings"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type FormBuilder struct {
	submitEndpoint string
	fields         []Renderer
}

type RawStringRenderer struct {
	raw string
}

func (r RawStringRenderer) Render(w io.Writer) {
	w.Write([]byte(r.raw))
}

func (b *FormBuilder) AddRaw(fieldName string, s string) *FormBuilder {
	var sb strings.Builder
	node := Div(Div(Class(FormFieldClass), Text(fieldName)), Raw(s))
	if err := node.Render(&sb); err != nil {
		panic(err)
	}
	b.fields = append(b.fields, RawStringRenderer{sb.String()})
	return b
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
