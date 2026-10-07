package userstore

import (
	"fmt"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type UserRepo interface {
	Find(u.UserId) (*u.BaseUser, error)
	List(pag.PaginationParams) ([]u.BaseUser, error)
	Count() (*int64, error)
}

type RoleRepo interface {
	Find(r.RoleName) (*r.Role, error)
}

type GroupRepo interface {
	Find(string) (*g.Group, error)
}

type UserStore struct {
	UserRepo
	GroupRepo
	RoleRepo
}

func (s UserStore) Find(id u.UserId) (*u.User, error) {
	errCtx := fmt.Errorf("fetching user with id %v", id)
	base, err := s.UserRepo.Find(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}

	var roles []r.Role
	for _, roleName := range base.Roles {
		role, err := s.RoleRepo.Find(roleName)
		if err != nil {
			return nil, fmt.Errorf("%w: aggregating role with name %v: %w", errCtx, roleName, err)
		}
		roles = append(roles, *role)
	}

	var groups []g.Group
	for _, groupName := range base.Groups {
		group, err := s.GroupRepo.Find(groupName)
		if err != nil {
			return nil, fmt.Errorf("%w: aggregating group with name %v: %w", errCtx, groupName, err)
		}
		groups = append(groups, *group)
	}

	user := u.NewUser(base.Id, u.WithRoles(roles), u.WithGroups(groups),
		u.WithPasswordHash(base.HashPassword), u.WithHashedPersonalAccessToken(base.HashPAT))
	return &user, nil
}

func (s UserStore) List(p pag.PaginationParams) ([]u.User, *pag.Pagination, error) {
	count, err := s.UserRepo.Count()
	if err != nil {
		return nil, nil, err
	}
	baseUsers, err := s.UserRepo.List(p)
	if err != nil {
		return nil, nil, err
	}
	var users []u.User
	for _, b := range baseUsers {
		user, err := s.Find(b.Id)
		if err != nil {
			return nil, nil, err
		}
		users = append(users, *user)
	}
	pagination := pag.New(p.Page, p.PageSize, *count)
	return users, &pagination, nil
}

func NewUserStore(ur UserRepo, rr RoleRepo, gr GroupRepo) UserStore {
	return UserStore{ur, gr, rr}
}
