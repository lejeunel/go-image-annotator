package user

import (
	"context"
	"slices"
	"time"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	r "github.com/lejeunel/go-image-annotator/entities/role"
)

var UserContextKey = "user"

type UserId = string

type ForgotPasswordState struct {
	Id        UserId
	ExpiresAt *time.Time
}

type BaseUser struct {
	Id           string
	HashPAT      []byte
	HashPassword []byte
	Groups       []string
	Roles        []string
}

func (u BaseUser) IsAdmin() bool {
	return slices.Contains(u.Roles, "admin")
}

type User struct {
	Id           string
	HashPAT      []byte
	HashPassword []byte
	Roles        []r.Role
	Groups       []g.Group
}

func (u User) ToBase() BaseUser {
	return BaseUser{
		Id: u.Id, HashPAT: u.HashPAT, HashPassword: u.HashPassword,
		Roles: u.RoleNames(), Groups: u.GroupNames(),
	}
}

func (u User) HasRole(role r.RoleName) bool {
	for _, r := range u.Roles {
		if r.Name == role {
			return true
		}
	}
	return false
}

func (u User) IsInGroup(group string) bool {
	for _, g := range u.Groups {
		if g.Name == group {
			return true
		}
	}
	return false
}

func (u User) IsAdmin() bool {
	for _, r := range u.Roles {
		if r.Name == "admin" {
			return true
		}
	}
	return false
}

func (u User) RoleNames() []string {
	var roleNames []string
	for _, r := range u.Roles {
		roleNames = append(roleNames, r.Name)
	}
	return roleNames
}

func (u User) GroupNames() []string {
	var groupNames []string
	for _, g := range u.Groups {
		groupNames = append(groupNames, g.Name)
	}
	return groupNames
}

func (u User) Permissions() []string {
	var permissions []string
	unqs := make(map[string]bool)
	for _, r := range u.Roles {
		for _, p := range r.Methods {
			if _, ok := unqs[p]; !ok {
				unqs[p] = true
				permissions = append(permissions, p)
			}
		}
	}
	return permissions
}

func NewUser(id UserId, opts ...Option) User {
	l := &User{Id: id}
	for _, opt := range opts {
		opt(l)
	}
	return *l
}

type Option func(*User)

func WithHashedPersonalAccessToken(h []byte) Option {
	return func(l *User) {
		l.HashPAT = h
	}
}

func WithPasswordHash(h []byte) Option {
	return func(l *User) {
		l.HashPassword = h
	}
}

func WithGroups(groups []g.Group) Option {
	return func(l *User) {
		l.Groups = groups
	}
}

func WithRoles(roles []r.Role) Option {
	return func(l *User) {
		l.Roles = roles
	}
}

func AppendUserToContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, UserContextKey, &user)
}

func IdentityFromContext(ctx context.Context) *User {
	v := ctx.Value(UserContextKey)
	if v == nil {
		return nil
	}
	user, ok := v.(*User)
	if !ok {
		return nil
	}

	return user
}
