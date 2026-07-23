package form

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed templates/selectable_combobox.html
var SelectableComboboxTemplate string

type SelectableComboboxField struct {
	Value    string
	Selected bool
}

type SelectableComboboxData struct {
	ID      string
	Title   string
	Options []SelectableComboboxField
}

type SelectableCombobox struct {
	title  string
	id     string
	fields []SelectableComboboxField
}

func NewSelectableCombobox(title, id string) SelectableCombobox {
	return SelectableCombobox{title: title, id: id}
}

func (f *SelectableCombobox) AddField(value string, selected bool) *SelectableCombobox {
	f.fields = append(f.fields, SelectableComboboxField{Value: value, Selected: selected})
	return f
}

func (f *SelectableCombobox) Render(w io.Writer) {
	t := template.New("")
	t.Parse(SelectableComboboxTemplate)
	if err := t.ExecuteTemplate(w, "selectable_combobox",
		SelectableComboboxData{ID: f.id, Title: f.title, Options: f.fields}); err != nil {
		panic(err)
	}

}
