package fake

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type GroupRepo struct {
	ErrOnGetGroupOfCollection error
	ErrOnExists               error
	ErrOnFind                 error
	ErrOnList                 error
	ErrOnUpdate               error
	Existing                  []grp.Group
	GotUpdate                 grp.UpdateModel
}

func (r *GroupRepo) Find(name string) (*grp.Group, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}
	for _, g := range r.Existing {
		if g.Name == name {
			return &g, nil
		}
	}

	return nil, e.ErrNotFound
}

func (r *GroupRepo) Exists(name string) (*bool, error) {
	if r.ErrOnExists != nil {
		return nil, r.ErrOnExists
	}
	exists := false
	for _, g := range r.Existing {
		if g.Name == name {
			exists = true
			return &exists, nil
		}
	}
	return &exists, nil
}

func (r *GroupRepo) List() ([]grp.Group, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	return r.Existing, nil
}

func (r *GroupRepo) Update(m grp.UpdateModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdate = m
	return nil
}
