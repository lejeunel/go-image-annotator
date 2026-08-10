package form

import (
	"io"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type InputNodeFunc func(string, *string, bool) Node

type FormField struct {
	fieldName   string
	displayName string
	value       *string
	required    bool
	InputNodeFunc
}
type FormFieldOption func(*FormField)

func WithRequired() FormFieldOption {
	return func(c *FormField) {
		c.required = true
	}
}

func WithDefault(value string) FormFieldOption {
	return func(c *FormField) {
		c.value = &value
	}
}

func NewFormTextField(fieldName, displayName string, opts ...FormFieldOption) FormField {
	f := &FormField{
		fieldName: fieldName, displayName: displayName,
		InputNodeFunc: makeInputNodeFunc("text", nil),
	}
	for _, opt := range opts {
		opt(f)
	}
	return *f
}

func NewFormIntField(fieldName, displayName string, opts ...FormFieldOption) FormField {
	f := &FormField{
		fieldName: fieldName, displayName: displayName,
		InputNodeFunc: makeInputNodeFunc("number", nil),
	}
	for _, opt := range opts {
		opt(f)
	}
	return *f
}

func NewFormFloatField(fieldName, displayName string, opts ...FormFieldOption) FormField {
	f := &FormField{
		fieldName: fieldName, displayName: displayName,
		InputNodeFunc: makeInputNodeFunc("number", Step("any")),
	}
	for _, opt := range opts {
		opt(f)
	}
	return *f
}

func NewFormDateTimeField(fieldName, displayName string, opts ...FormFieldOption) FormField {
	f := &FormField{
		fieldName: fieldName, displayName: displayName,
		InputNodeFunc: makeInputNodeFunc("datetime-local", nil),
	}
	for _, opt := range opts {
		opt(f)
	}
	return *f
}

func NewFormBooleanField(fieldName, displayName string, opts ...FormFieldOption) FormField {
	f := &FormField{
		fieldName: fieldName, displayName: displayName,
		InputNodeFunc: makeInputNodeFunc("checkbox", nil),
	}
	for _, opt := range opts {
		opt(f)
	}
	return *f
}

func (f FormField) label() Node {
	displayName := f.displayName
	return Label(
		For(f.fieldName),
		Text(displayName),
		Class("w-fit pl-0.5 text-sm text-on-surface dark:text-on-surface-dark"),
	)
}

func makeInputNodeFunc(inputType string, extra ...Node) InputNodeFunc {
	return func(fieldName string, defaultValue *string, required bool) Node {
		var value string
		if defaultValue != nil {
			value = *defaultValue
		}

		nodes := []Node{
			Type(inputType),
			ID(fieldName),
			Name(fieldName),
		}

		nodes = append(nodes, extra...)

		nodes = append(nodes,
			If(required, Required()),
			If(value != "", Value(value)),
			Class(
				"w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent",
			),
		)

		return Input(nodes...)
	}
}

func (f FormField) input() Node {
	return f.InputNodeFunc(f.fieldName, f.value, f.required)
}

func (f FormField) Render(w io.Writer) {
	Div(Class("w-full max-w-xs flex flex-col gap-1"),
		Group([]Node{Div(Class("mr-2"), f.label()), f.input()})).Render(w)
}
