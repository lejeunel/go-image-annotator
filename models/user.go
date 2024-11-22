package models

import (
	"context"
	"errors"
	"strings"
)

type User struct {
	Email string
	Roles []string
}

func GetUserFromContext(ctx context.Context) (*User, error) {
	email := ctx.Value("user_email")
	if email == nil {
		return nil, errors.New("could not find user_email in context")
	}

	roles := ctx.Value("user_roles")
	if roles == nil {
		return nil, errors.New("could not find user_roles in context")
	}

	roles_arr := strings.Split(roles.(string), ",")
	user := User{Email: email.(string), Roles: roles_arr}

	return &user, nil

}
