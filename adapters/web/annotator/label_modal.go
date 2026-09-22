package annotator

import (
	"bytes"
	_ "embed"
	"text/template"

	s "github.com/lejeunel/go-image-annotator/adapters/web/components/select"
)

//go:embed templates/label_modal.html
var LabelModal string

func NewLabelModal(labels []string) string {
	tModal := template.New("")
	template.Must(tModal.Parse(LabelModal))
	template.Must(tModal.Parse(s.SingleSearch))

	var buf bytes.Buffer
	if err := tModal.ExecuteTemplate(
		&buf,
		"label_modal",
		s.SearchableData{Items: labels, FieldName: "label"},
	); err != nil {
		panic(err)
	}

	return buf.String()
}
