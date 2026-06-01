package annotator

import (
	"bytes"
	"text/template"
)

type NewLabelModal struct {
	Labels         []string
	SelectorIsOpen bool
	Selected       *string
	SubmissionFn   string
	ModalName      string
}

type LabelModalKind int

const (
	RegionLabelModal LabelModalKind = iota
	ImageLabelModal
)

func makeLabelModal(labels []string, kind LabelModalKind) (*string, error) {
	tModal := template.New("")
	template.Must(tModal.ParseFS(templatesFiles, "templates/label_modal_combobox.html"))
	template.Must(tModal.ParseFS(templatesFiles, "templates/label_modal.html"))

	var modal NewLabelModal
	switch kind {
	case RegionLabelModal:
		modal = NewLabelModal{labels, true, nil, "submit_region", "regionLabelModal"}
	default:
		modal = NewLabelModal{labels, true, nil, "submit_label", "imageLabelModal"}
	}

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(&buf, "label_modal",
		modal); err != nil {
		return nil, err
	}

	str := buf.String()
	return &str, nil
}
