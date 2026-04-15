package annotator

import (
	"bytes"
	"text/template"
)

func makeLabelModal(labels []string) (string, error) {
	tLabelModal, err := template.New("labelModal").ParseFS(templatesFiles,
		"templates/label_selector.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := struct {
		Labels []string
	}{labels}
	err = tLabelModal.ExecuteTemplate(&buf, "labelModal",
		data)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
