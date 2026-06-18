package sqlite

import (
	"github.com/lejeunel/go-image-annotator/modules/auth"
	tok "github.com/lejeunel/go-image-annotator/modules/token"
	usr "github.com/lejeunel/go-image-annotator/use-cases/user"
	agr "github.com/lejeunel/go-image-annotator/use-cases/user/assign-group"
	aro "github.com/lejeunel/go-image-annotator/use-cases/user/assign-role"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	"github.com/lejeunel/go-image-annotator/use-cases/user/read"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	adm "github.com/lejeunel/go-image-annotator/use-cases/user/set-admin"
	ugr "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-group"
	uar "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-role"
)

func NewSQLiteUserInteractors(repos SQLiteRepos, tokenGen tok.TokenGenerator,
	auth auth.Auth) usr.Interactors {
	return usr.Interactors{
		Find:          read.New(repos.User, read.WithAuth(auth)),
		Create:        create.New(repos.User, tokenGen, create.WithAuth(auth)),
		Delete:        delete.New(repos.User, delete.WithAuth(auth)),
		List:          list.New(repos.User, list.WithAuth(auth)),
		RenewToken:    rt.New(repos.User, tokenGen, rt.WithAuth(auth)),
		AssignRole:    aro.New(repos.User, aro.WithAuth(auth)),
		UnAssignRole:  uar.New(repos.User, uar.WithAuth(auth)),
		AssignGroup:   agr.New(repos.User, repos.Group, agr.WithAuth(auth)),
		UnAssignGroup: ugr.New(repos.User, repos.Group, ugr.WithAuth(auth)),
		SetAdmin:      adm.New(repos.User, adm.WithAuth(auth)),
	}
}
