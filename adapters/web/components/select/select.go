package selector

import (
	"bytes"
	_ "embed"
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

type MultiselectLabelData struct {
	Labels    []string
	FieldName string
}

type LabelModalKind int

const (
	RegionLabelModal LabelModalKind = iota
	ImageLabelModal
)

func NewSingleSelectCombobox(items []string, fieldName string) string {
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

func NewMultiSelectCombobox(labels []string, fieldName string) string {
	tModal := template.New("")
	template.Must(tModal.Parse(MultiSearch))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(
		&buf,
		"multi_search",
		MultiselectLabelData{Labels: labels, FieldName: fieldName},
	); err != nil {
		panic(err)
	}

	return buf.String()
}
