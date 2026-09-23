package selector

import (
	"bytes"
	_ "embed"
	"io"
	"text/template"
)

//go:embed single_search.html
var SingleSearch string

//go:embed multi_search.html
var MultiSearch string

type SearchableData struct {
	Items     []string
	FieldName string
}

type Item struct {
	Name     string
	IsActive bool
}

type MultiSelectData struct {
	Items     []Item
	FieldName string
}

type LabelModalKind int

const (
	RegionLabelModal LabelModalKind = iota
	ImageLabelModal
)

func NewSingleSelect(items []string, fieldName string) string {
	tModal := template.New("")
	template.Must(tModal.Parse(SingleSearch))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(
		&buf,
		"search_combobox",
		SearchableData{Items: items, FieldName: fieldName},
	); err != nil {
		panic(err)
	}

	return buf.String()
}

type MultiSelectBuilder struct {
	fieldName string
	items     []Item
}

func NewMultiSelectBuilder(fieldName string) MultiSelectBuilder {
	return MultiSelectBuilder{fieldName: fieldName}
}

func (b *MultiSelectBuilder) AddItem(name string, isActive bool) *MultiSelectBuilder {
	b.items = append(b.items, Item{Name: name, IsActive: isActive})
	return b
}

func (b MultiSelectBuilder) Render(w io.Writer) {
	tModal := template.New("")
	template.Must(tModal.Parse(MultiSearch))

	if err := tModal.ExecuteTemplate(
		w,
		"multi_search",
		MultiSelectData{Items: b.items, FieldName: b.fieldName},
	); err != nil {
		panic(err)
	}
}
