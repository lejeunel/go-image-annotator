package fake

import (
	"slices"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type RoleRepo struct {
	ErrOnCreate   error
	ErrOnDelete   error
	ErrOnFind     error
	ErrOnList     error
	ErrOnExists   error
	ErrOnUpdate   error
	ExistingNames []string
	ExistingRoles []rl.Role
	Created       []rl.Role
	IsAssigned_   bool
	Return        rl.Role
	ReturnList    []rl.Role
	GotUpdatable  rl.UpdatableModel
}

func (r *RoleRepo) Create(role rl.Role) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}

	r.Created = append(r.Created, role)
	return nil
}

func (r *RoleRepo) Exists(name string) (*bool, error) {
	if r.ErrOnExists != nil {
		return nil, r.ErrOnExists
	}

	exist := true
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
	for _, r := range r.ExistingRoles {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, e.ErrNotFound
}

func (r *RoleRepo) List() ([]rl.Role, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	return r.ReturnList, nil
}

func (r *RoleRepo) Update(m rl.UpdatableModel) error {
	if r.ErrOnUpdate != nil {
		return r.ErrOnUpdate
	}
	r.GotUpdatable = m
	return nil
}
