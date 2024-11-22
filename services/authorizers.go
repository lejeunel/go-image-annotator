package services

import (
	"context"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"slices"
)

func CheckAuthorization(ctx context.Context, authorizedRole string) error {
	user, err := m.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	if !(slices.Contains(user.Roles, "admin") || slices.Contains(user.Roles, authorizedRole)) {
		return e.ErrPermission{Operation: "Contribute image", NeededRole: authorizedRole, UserRoles: user.Roles}
	}

	return nil

}
