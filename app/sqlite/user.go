package sqlite

import (
	sqlitegrp "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	sqliterol "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	sqliteusr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
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
	userRepo sqliteusr.SQLiteUserRepo,
	grpRepo sqlitegrp.SQLiteGroupRepo,
	roleRepo sqliterol.SQLiteRoleRepo,
	ApitokenGen create.APITokenGenerator,
	forgotPasswordTokenGen fp.TokenGenerator,
	passwordValidator pw.PasswordValidator,
	passwordHasher tk.TokenHasher,
	passwordVerifier cpw.TokenVerifier,
	forgotPassworkTokenExpirationMinutes int,
	pwGen create.PasswordGenerator,
	auth auth.Interface,
) usr.Interactors {
	return usr.Interactors{
		Find:             find.New(userRepo, find.WithAuth(auth)),
		Create:           create.New(userRepo, ApitokenGen, pwGen, create.WithAuth(auth)),
		Delete:           delete.New(userRepo, delete.WithAuth(auth)),
		List:             list.New(userRepo, list.WithAuth(auth)),
		RenewToken:       rt.New(userRepo, ApitokenGen),
		UpdatePrivileges: upr.New(userRepo, grpRepo, roleRepo, upr.WithAuth(auth)),
		RequestForgottenPassword: fp.New(
			userRepo,
			forgotPassworkTokenExpirationMinutes,
			forgotPasswordTokenGen,
		),
		ResetForgottenPassword: rfpw.New(userRepo, passwordHasher, passwordValidator),
		ChangePassword:         cpw.New(userRepo, passwordVerifier, passwordValidator),
	}
}
