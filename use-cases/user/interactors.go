package user

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	cpw "github.com/lejeunel/go-image-annotator/use-cases/user/change-password"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/user/find"
	fp "github.com/lejeunel/go-image-annotator/use-cases/user/forgot-password"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	rat "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	rfpw "github.com/lejeunel/go-image-annotator/use-cases/user/reset-forgotten-password"
	upr "github.com/lejeunel/go-image-annotator/use-cases/user/update-privileges"
)

type Interactors struct {
	Find                     find.Interactor
	Create                   create.Interactor
	Delete                   delete.Interactor
	List                     list.Interactor
	RenewToken               rat.Interactor
	UpdatePrivileges         upr.Interactor
	RequestForgottenPassword fp.Interactor
	ResetForgottenPassword   rfpw.Interactor
	ChangePassword           cpw.Interactor
	DefaultPageSize          int
	Authorizer               auth.Authorizer
}
