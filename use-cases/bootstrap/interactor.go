package bootstrap

import (
	"bytes"
	"context"
	"fmt"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	a "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
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
	FileStore fs.FileStore
}

func New(ur UserRepo, rr RoleRepo, f fs.FileStore,
	h PasswordHasher, v pw.PasswordValidator,
) Interactor {
	return Interactor{ur, rr, h, v, f}
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

	for _, defaultRole := range rl.DefaultRoleNames {
		if err := i.RoleRepo.Create(
			rl.NewRole(
				rl.NewRoleId(),
				defaultRole.Name,
				rl.WithDescription(defaultRole.Description))); err != nil {
			out.Error(
				fmt.Errorf(
					"%w: creating role %v: %v: %w",
					errCtx,
					defaultRole.Name,
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

	user := u.NewUser(
		r.InitialAdminEmail,
		u.WithPasswordHash(pwHash),
		u.WithRoles([]string{"admin"}),
	)
	if err := i.UserRepo.Create(user); err != nil {
		out.Error(fmt.Errorf("%w: creating admin user: %v: %w", errCtx, err, e.ErrInternal))
		return
	}

	var buf bytes.Buffer
	if err := a.MarshalPolicies(a.DefaultPolicies, &buf); err != nil {
		out.Error(
			fmt.Errorf("%w: generating default yaml policies: %v: %w", errCtx, err, e.ErrInternal),
		)
		return
	}
	if err := i.FileStore.Store(a.DefaultPolicyFileName, &buf); err != nil {
		out.Error(
			fmt.Errorf("%w: writing default yaml policies: %v: %w", errCtx, err, e.ErrInternal),
		)
		return
	}

	out.SuccessBootstrap(Response{Skipped: false})
}
