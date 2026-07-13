package role

import (
	ro "github.com/lejeunel/go-image-annotator/entities/role"
)

func CreateRole(repo SQLiteRoleRepo, name string) (*ro.Role, error) {
	r := ro.NewRole(ro.NewRoleId(), name,
		ro.WithDescription("a-description"))

	if err := repo.Create(r); err != nil {
		return nil, err
	}
	return &r, nil

}
