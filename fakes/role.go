package fake

import (
	"slices"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type RoleRepo struct {
	ErrOnCreate   error
	ErrOnDelete   error
	ErrOnFind     error
	ErrOnCount    error
	ErrOnList     error
	ErrOnUpdate   error
	ExistingNames []string
	Got           *rl.Role
	IsAssigned_   bool
	Count_        int64
	Return        rl.Role
	GotUpdatable  rl.UpdatableModel
}

func (r *RoleRepo) Create(role rl.Role) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Got = &role
	return nil
}

func (r *RoleRepo) Exists(name string) (*bool, error) {
	var exist = true
	if slices.Contains(r.ExistingNames, name) {
		return &exist, nil
	}
	exist = false
	return &exist, nil
}

func (r *RoleRepo) Delete(string) error {

	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}
	return nil
}

func (r *RoleRepo) IsAssigned(c string) (*bool, error) {
	res := true
	if r.IsAssigned_ {
		return &res, nil
	}
	res = false
	return &res, nil
}

func (r *RoleRepo) Find(name string) (*rl.Role, error) {
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}

	return &r.Return, nil

}

func (r *RoleRepo) Count() (*int64, error) {
	count := int64(0)
	if r.ErrOnCount != nil {
		return &count, r.ErrOnCount
	}
	return &r.Count_, nil
}

func (r *RoleRepo) List(req pag.PaginationParams) ([]*rl.Role, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	result := []*rl.Role{}
	for range req.PageSize {
		c := rl.NewRole(rl.NewRoleId(), "a-role")
		result = append(result, &c)
	}
	return result, nil
}

func (r *RoleRepo) Update(m rl.UpdatableModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdatable = m
	return nil
}
