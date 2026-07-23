package sqlite

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	pw "github.com/lejeunel/go-image-annotator/modules/password-validator"
	tk "github.com/lejeunel/go-image-annotator/modules/token"
	usr "github.com/lejeunel/go-image-annotator/use-cases/user"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/user/find"
	fp "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	rfpw "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
	upr "github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

func NewSQLiteUserInteractors(
	repos SQLiteRepos,
	ApitokenGen create.APITokenGenerator,
	forgotPasswordTokenGen fp.TokenGenerator,
	passwordValidator pw.PasswordValidator,
	passwordHasher tk.TokenHasher,
	passwordVerifier cpw.TokenVerifier,
	forgotPassworkTokenExpirationMinutes int,
	pwGen create.PasswordGenerator,
	auth auth.Authorizer) usr.Interactors {
	return usr.Interactors{
		Find:                     find.New(repos.User, find.WithAuth(auth)),
		Create:                   create.New(repos.User, ApitokenGen, pwGen, create.WithAuth(auth)),
		Delete:                   delete.New(repos.User, delete.WithAuth(auth)),
		List:                     list.New(repos.User, list.WithAuth(auth)),
		RenewToken:               rt.New(repos.User, ApitokenGen, rt.WithAuth(auth)),
		UpdatePrivileges:         upr.New(repos.User, repos.Group, repos.Role, upr.WithAuth(auth)),
		RequestForgottenPassword: fp.New(repos.User, forgotPassworkTokenExpirationMinutes, forgotPasswordTokenGen, fp.WithAuth(auth)),
		ResetForgottenPassword:   rfpw.New(repos.User, passwordHasher, passwordValidator),
		ChangePassword:           cpw.New(repos.User, passwordVerifier, passwordValidator, cpw.WithAuth(auth)),
	}
}
