package fake

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"slices"
)

type GroupRepo struct {
	MissingGroup              bool
	ErrOnGetGroupOfCollection error
	ErrOnExists               error
	Return                    string
	ExistingNames             []string
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

func (r *GroupRepo) Exists(name string) (*bool, error) {
	if r.ErrOnExists != nil {
		return nil, r.ErrOnExists
	}
	var res bool
	if slices.Contains(r.ExistingNames, name) {
		res = true
		return &res, nil
	}
	res = false
	return &res, e.ErrNotFound
}
