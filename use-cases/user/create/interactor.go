package create

import (
	"context"
	"fmt"

	tk "github.com/lejeunel/go-image-annotator/entities/token"
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type APITokenGenerator interface {
	Generate() (*tk.Token, error)
}

type PasswordGenerator interface {
	Generate() (*tk.Token, error)
	Hash(string) []byte
}

type Interactor struct {
	Repo              Repo
	GroupRepo         GroupRepo
	RoleRepo          RoleRepo
	APITokenGenerator APITokenGenerator
	PasswordGenerator PasswordGenerator

	auth Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "creating user"
	if err := i.auth.CreateUser(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return

	}
	if err := i.checkDuplicate(r.Id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.checkGroupsExist(r.Groups); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.checkRolesExist(r.Roles); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	token, err := i.APITokenGenerator.Generate()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	var passwordHash []byte
	if r.Password != nil {
		passwordHash = i.PasswordGenerator.Hash(*r.Password)
	} else {
		passwordPair, err := i.PasswordGenerator.Generate()
		passwordHash = passwordPair.Hash
		if err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}
	user := usr.BaseUser{
		Id:           r.Id,
		HashPAT:      token.Hash,
		HashPassword: passwordHash,
	}
	if r.Groups != nil {
		user.Groups = *r.Groups
	}
	if r.Roles != nil {
		user.Roles = *r.Roles
	}
	if err := i.Repo.Create(user); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessCreateUser(Response{
		Id:     user.Id,
		Groups: user.Groups,
		Roles:  user.Roles,
	})
}

func (i *Interactor) checkDuplicate(id string) error {
	errBaseMsg := "checking for duplicate user with id %v: %w"
	alreadyExists, err := i.Repo.Exists(id)
	if err != nil {
		return fmt.Errorf(errBaseMsg, id, e.ErrInternal)
	}
	if alreadyExists {
		return fmt.Errorf(errBaseMsg, id, e.ErrDuplicate)
	}
	return nil
}

func (i *Interactor) checkGroupsExist(groups *[]string) error {
	if groups == nil {
		return nil
	}
	for _, g := range *groups {
		exists, err := i.GroupRepo.Exists(g)
		if err != nil {
			return fmt.Errorf("checking whether group %v exists: %w", g, err)
		}
		if !*exists {
			return fmt.Errorf("checking whether group %v exists: %w", g, e.ErrNotFound)
		}
	}
	return nil
}

func (i *Interactor) checkRolesExist(roles *[]string) error {
	if roles == nil {
		return nil
	}
	for _, role := range *roles {
		exists, err := i.RoleRepo.Exists(role)
		if err != nil {
			return fmt.Errorf("checking whether role %v exists: %w", role, err)
		}
		if !*exists {
			return fmt.Errorf("checking whether role %v exists: %w", role, e.ErrNotFound)
		}
	}
	return nil
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo,
	gr GroupRepo,
	rr RoleRepo,
	tg APITokenGenerator,
	pg PasswordGenerator, opts ...Option,
) Interactor {
	i := &Interactor{
		Repo:              r,
		GroupRepo:         gr,
		RoleRepo:          rr,
		auth:              auth.NewVoidAuth(),
		APITokenGenerator: tg,
		PasswordGenerator: pg,
	}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
