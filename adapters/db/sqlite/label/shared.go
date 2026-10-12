package label

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

func CreateLabel(repo LabelRepo, name string) (*lbl.Label, error) {
	label := lbl.NewLabel(lbl.NewLabelId(), name, lbl.WithDescription("a-description"))
	if err := repo.Create(label); err != nil {
		return nil, err
	}
	return &label, nil
}
