package fake

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type GroupRepo struct {
	MissingGroup              bool
	ErrOnGetGroupOfCollection error
	Return                    string
}

func (r *GroupRepo) Find(name string) (*grp.Group, error) {
	if r.MissingGroup {
		return nil, e.ErrNotFound
	}

	return &grp.Group{Name: r.Return}, nil
}

func (r *GroupRepo) GroupOfCollection(string) (*string, error) {
	if r.ErrOnGetGroupOfCollection != nil {
		return nil, r.ErrOnGetGroupOfCollection
	}
	return &r.Return, nil
}
