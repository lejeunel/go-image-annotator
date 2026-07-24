package form

import (
	"io"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type FormTextField struct {
	fieldName   string
	displayName string
	divId       string
	value       *string
	required    bool
}
type FormTextFieldOption func(*FormTextField)

func WithRequired() FormTextFieldOption {
	return func(c *FormTextField) {
		c.required = true
	}
}
func WithDefault(value string) FormTextFieldOption {
	return func(c *FormTextField) {
		c.value = &value
	}
}

func NewFormTextField(fieldName, displayName, divId string, opts ...FormTextFieldOption) FormTextField {
	f := &FormTextField{fieldName: fieldName, displayName: displayName, divId: divId}
	for _, opt := range opts {
		opt(f)
	}
	return *f
}

func (f FormTextField) label() Node {
	displayName := f.displayName
	if !f.required {
		displayName = displayName + " (Optional)"
	}
	return Label(For(f.fieldName), Text(displayName), Class("block text-sm font-medium text-gray-900 dark:text-white"))
}

func (f FormTextField) input() Node {
	var value string
	if f.value != nil {
		value = *f.value
	}
	return Input(Type("text"), ID(f.divId), Name(f.fieldName), If(f.required, Required()), Value(value),
		Class("w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"))
}

func (f FormTextField) Render(w io.Writer) {
	Group([]Node{f.label(), f.input()}).Render(w)

}
