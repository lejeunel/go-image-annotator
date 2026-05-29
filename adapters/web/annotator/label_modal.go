package annotator

import (
	"bytes"
	"text/template"
)

type NewLabelSelector struct {
	Labels         []string
	SelectorIsOpen bool
	Selected       *string
}

func makeLabelModal(labels []string) (*string, error) {
	tModal := template.New("")
	template.Must(tModal.ParseFS(templatesFiles, "templates/label_modal_combobox.html"))
	template.Must(tModal.ParseFS(templatesFiles, "templates/label_modal.html"))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(&buf, "label_modal",
		NewLabelSelector{labels, true, nil}); err != nil {
		return nil, err
	}

	str := buf.String()
	return &str, nil
}
