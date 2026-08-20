package form

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed templates/combobox.html
var ComboboxTemplate string

type ComboboxField struct {
	Value string
}

type ComboboxData struct {
	ID            string
	Title         string
	Value         string
	Options       []ComboboxField
	SelectedValue string
	DefaultValue  string
}

type Combobox struct {
	title         string
	id            string
	fields        []ComboboxField
	selectedValue *string
	defaultValue  string
}

func NewCombobox(title, id, d string) Combobox {
	return Combobox{title: title, id: id, defaultValue: d}
}

func (f *Combobox) AddField(value string) *Combobox {
	f.fields = append(f.fields, ComboboxField{Value: value})
	return f
}

func (f *Combobox) SetSelectedValue(value string) *Combobox {
	f.selectedValue = &value
	return f
}

func (f *Combobox) Render(w io.Writer) {
	t := template.New("")
	_, err := t.Parse(ComboboxTemplate)
	if err != nil {
		panic(err)
	}
	data := ComboboxData{ID: f.id, Title: f.title, Options: f.fields, DefaultValue: f.defaultValue}
	if f.selectedValue != nil {
		data.SelectedValue = *f.selectedValue
	} else {
		data.SelectedValue = "null"
	}
	if err := t.ExecuteTemplate(w, "combobox", data); err != nil {
		panic(err)
	}
}
