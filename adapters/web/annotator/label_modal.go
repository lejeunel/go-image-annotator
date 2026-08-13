package annotator

import (
	"bytes"
	"text/template"
)

type NewLabelModal struct {
	Labels         []string
	SelectorIsOpen bool
	Selected       *string
}

type LabelModalKind int

const (
	RegionLabelModal LabelModalKind = iota
	ImageLabelModal
)

func makeLabelModal(labels []string) string {
	tModal := template.New("")
	template.Must(tModal.ParseFS(templatesFiles, "templates/label_modal_search_combobox.html"))
	template.Must(tModal.ParseFS(templatesFiles, "templates/label_modal.html"))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(
		&buf,
		"label_modal",
		NewLabelModal{Labels: labels},
	); err != nil {
		panic(err)
	}

	return buf.String()
}
