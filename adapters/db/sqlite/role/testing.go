package role

import (
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"

	ro "github.com/lejeunel/go-image-annotator/entities/role"
)

func CreateRole(repo RoleRepo, name string) (*ro.Role, error) {
	r := ro.NewRole(ro.NewRoleId(), name,
		ro.WithDescription("a-description"))

	if err := repo.Create(r); err != nil {
		return nil, err
	}
	return &r, nil
}

func NewTestRoleRepo() RoleRepo {
	return NewRoleRepo(s.NewInMemory())
}
