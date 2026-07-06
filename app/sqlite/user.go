package sqlite

import (
	"github.com/lejeunel/go-image-annotator/modules/auth"
	usr "github.com/lejeunel/go-image-annotator/use-cases/user"
	agr "github.com/lejeunel/go-image-annotator/use-cases/user/assign-group"
	aro "github.com/lejeunel/go-image-annotator/use-cases/user/assign-role"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	fp "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	"github.com/lejeunel/go-image-annotator/use-cases/user/read"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	rfp "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
	adm "github.com/lejeunel/go-image-annotator/use-cases/user/set-admin"
	ugr "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-group"
	uar "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-role"
)

func NewSQLiteUserInteractors(
	repos SQLiteRepos,
	ApitokenGen create.APITokenGenerator,
	forgotPasswordTokenGen fp.TokenGenerator,
	passwordValidator rfp.PasswordValidator,
	passwordHasher rfp.TokenHasher,
	forgotPassworkTokenExpirationMinutes int,
	pwGen create.PasswordGenerator,
	auth auth.Auth) usr.Interactors {
	return usr.Interactors{
		Find:                     read.New(repos.User, read.WithAuth(auth)),
		Create:                   create.New(repos.User, ApitokenGen, pwGen, create.WithAuth(auth)),
		Delete:                   delete.New(repos.User, delete.WithAuth(auth)),
		List:                     list.New(repos.User, list.WithAuth(auth)),
		RenewToken:               rt.New(repos.User, ApitokenGen, rt.WithAuth(auth)),
		AssignRole:               aro.New(repos.User, aro.WithAuth(auth)),
		UnAssignRole:             uar.New(repos.User, uar.WithAuth(auth)),
		AssignGroup:              agr.New(repos.User, repos.Group, agr.WithAuth(auth)),
		UnAssignGroup:            ugr.New(repos.User, repos.Group, ugr.WithAuth(auth)),
		SetAdmin:                 adm.New(repos.User, adm.WithAuth(auth)),
		RequestForgottenPassword: fp.New(repos.User, forgotPassworkTokenExpirationMinutes, forgotPasswordTokenGen, fp.WithAuth(auth)),
		ResetForgottenPassword:   rfp.New(repos.User, passwordHasher, passwordValidator),
	}
}
