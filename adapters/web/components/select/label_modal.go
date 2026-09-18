package templates

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed label_modal.html
var LabelModal string

//go:embed label_modal_search_combobox.html
var LabelModalSearchCombobox string

//go:embed multilabel_search.html
var MultilabelSearch string

type LabelModalData struct {
	Labels []string
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

func NewLabelModal(labels []string) string {
	tModal := template.New("")
	template.Must(tModal.Parse(LabelModalSearchCombobox))
	template.Must(tModal.Parse(LabelModal))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(
		&buf,
		"label_modal",
		LabelModalData{Labels: labels},
	); err != nil {
		panic(err)
	}

	return buf.String()
}

func NewMultiLabelCombobox(labels []string, fieldName string) string {
	tModal := template.New("")
	template.Must(tModal.Parse(MultilabelSearch))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(
		&buf,
		"multilabel_search",
		MultiselectLabelData{Labels: labels, FieldName: fieldName},
	); err != nil {
		panic(err)
	}

	return buf.String()
}
