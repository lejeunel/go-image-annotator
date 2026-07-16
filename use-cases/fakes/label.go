package fake

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type LabelRepo struct {
	Err         error
	Return      lbl.Label
	FetchedName string
}

func (r *LabelRepo) FindLabel(name string) (*lbl.Label, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.FetchedName = name
	return &r.Return, nil
}
