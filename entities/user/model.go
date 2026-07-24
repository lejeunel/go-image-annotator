package user

import (
	"context"
	"slices"
	"time"
)

var UserContextKey = "user"

type UserId = string

type ForgotPasswordState struct {
	Id        UserId
	ExpiresAt *time.Time
}

type User struct {
	Id           string
	HashPAT      []byte
	HashPassword []byte
	Roles        []string
	Groups       []string
}

func NewUser(id UserId, opts ...Option) User {
	l := &User{Id: id}
	for _, opt := range opts {
		opt(l)
	}
	return *l
}

func (u User) IsAdmin() bool {
	return slices.Contains(u.Roles, "admin")
}

type Option func(*User)

func WithHashedPersonalAccessToken(h []byte) Option {
	return func(l *User) {
		l.HashPAT = h
	}
}

func WithHashedPassword(h []byte) Option {
	return func(l *User) {
		l.HashPassword = h
	}
}

func WithGroups(groups []string) Option {
	return func(l *User) {
		l.Groups = groups
	}
}

func WithRoles(roles []string) Option {
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
