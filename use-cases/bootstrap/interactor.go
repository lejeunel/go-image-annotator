package bootstrap

import (
	"context"
	"fmt"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type PasswordHasher interface {
	Hash(string) []byte
}

type Interactor struct {
	UserRepo
	RoleRepo
	PasswordHasher
	pw.PasswordValidator
}

func New(ur UserRepo, rr RoleRepo,
	h PasswordHasher, v pw.PasswordValidator,
) Interactor {
	return Interactor{ur, rr, h, v}
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("bootstrapping application with initial admin user")
	adminExists, err := i.RoleRepo.Exists("admin")
	if err != nil {
		out.Error(
			fmt.Errorf("%w: checking existence of admin role: %v: %w", errCtx, err, e.ErrInternal),
		)
		return
	}
	if *adminExists {
		out.SuccessBootstrap(Response{Skipped: true})
		return
	}

	for _, policy := range a.DefaultPolicies {
		if err := i.RoleRepo.Create(
			rl.NewRole(
				rl.NewRoleId(),
				policy.Role,
				rl.WithDescription(policy.Description),
				rl.WithMethods(policy.Methods))); err != nil {
			out.Error(
				fmt.Errorf(
					"%w: creating role %v: %v: %w",
					errCtx,
					policy.Role,
					err,
					e.ErrInternal,
				),
			)
			return
		}
	}

	if err := i.PasswordValidator.Validate(r.InitialAdminPassword); err != nil {
		out.Error(fmt.Errorf("%w: validating initial password: %w", errCtx, err))
		return
	}

	pwHash := i.PasswordHasher.Hash(r.InitialAdminPassword)

	user := u.BaseUser{
		Id:           r.InitialAdminEmail,
		HashPassword: pwHash,
		Roles:        []string{"admin"},
	}
	if err := i.UserRepo.Create(user); err != nil {
		out.Error(fmt.Errorf("%w: creating admin user: %v: %w", errCtx, err, e.ErrInternal))
		return
	}

	out.SuccessBootstrap(Response{Skipped: false})
}
