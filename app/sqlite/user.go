package sqlite

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	usr "github.com/lejeunel/go-image-annotator/use-cases/user"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/user/find"
	fp "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	rfpw "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
	adm "github.com/lejeunel/go-image-annotator/use-cases/user/set-admin"
	ugr "github.com/lejeunel/go-image-annotator/use-cases/user/update-groups"
	uro "github.com/lejeunel/go-image-annotator/use-cases/user/update-roles"
)

func NewSQLiteUserInteractors(
	repos SQLiteRepos,
	ApitokenGen create.APITokenGenerator,
	forgotPasswordTokenGen fp.TokenGenerator,
	passwordValidator pw.PasswordValidator,
	passwordHasher tk.TokenHasher,
	forgotPassworkTokenExpirationMinutes int,
	pwGen create.PasswordGenerator,
	auth auth.Authorizer) usr.Interactors {
	return usr.Interactors{
		Find:                     find.New(repos.User, find.WithAuth(auth)),
		Create:                   create.New(repos.User, ApitokenGen, pwGen, create.WithAuth(auth)),
		Delete:                   delete.New(repos.User, delete.WithAuth(auth)),
		List:                     list.New(repos.User, list.WithAuth(auth)),
		RenewToken:               rt.New(repos.User, ApitokenGen, rt.WithAuth(auth)),
		UpdateRole:               uro.New(repos.User, uro.WithAuth(auth)),
		UpdateGroup:              ugr.New(repos.User, repos.Group, ugr.WithAuth(auth)),
		SetAdmin:                 adm.New(repos.User, adm.WithAuth(auth)),
		RequestForgottenPassword: fp.New(repos.User, forgotPassworkTokenExpirationMinutes, forgotPasswordTokenGen, fp.WithAuth(auth)),
		ResetForgottenPassword:   rfpw.New(repos.User, passwordHasher, passwordValidator),
	}
}
