package update

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Request struct {
	Name           string
	NewName        string
	NewDescription *string
	NewGroup       *string
	NewLabels      []lbl.LabelName
}
