package fake

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	"slices"
)

type GroupRepo struct {
	ErrOnGetGroupOfCollection error
	ErrOnExists               error
	ErrOnFind                 error
	ErrOnList                 error
	ErrOnUpdate               error
	Return                    grp.Group
	ReturnList                []grp.Group
	ExistingNames             []string
	GotUpdate                 grp.UpdateModel
}

func (r *GroupRepo) Find(name string) (*grp.Group, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}

	return &r.Return, nil
}

func (r *GroupRepo) GroupOfCollection(string) (*string, error) {
	if r.ErrOnGetGroupOfCollection != nil {
		return nil, r.ErrOnGetGroupOfCollection
	}
	return &r.Return.Name, nil
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
	return &res, nil
}

func (r *GroupRepo) List() ([]grp.Group, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	return r.ReturnList, nil
}

func (r *GroupRepo) Update(m grp.UpdateModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdate = m
	return nil
}
