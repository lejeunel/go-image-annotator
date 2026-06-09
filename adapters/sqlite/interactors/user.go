package interactors

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/sqlite"
	tg "github.com/lejeunel/go-image-annotator/app/token-generator"
	usr "github.com/lejeunel/go-image-annotator/use-cases/user"
	agr "github.com/lejeunel/go-image-annotator/use-cases/user/assign-group"
	aro "github.com/lejeunel/go-image-annotator/use-cases/user/assign-role"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	"github.com/lejeunel/go-image-annotator/use-cases/user/read"
	rt "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	ugr "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-group"
	uar "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-role"
)

func NewSQLiteUserInteractors(repos *infra.SQLiteInfra, tokenGen tg.TokenGenerator) *usr.Interactors {
	return &usr.Interactors{
		Find:          read.New(repos.User),
		Create:        create.New(repos.User, tokenGen),
		Delete:        delete.New(repos.User),
		List:          list.New(repos.User),
		RenewToken:    rt.New(repos.User, tokenGen),
		AssignRole:    aro.New(repos.User),
		UnAssignRole:  uar.New(repos.User),
		AssignGroup:   agr.New(repos.User, repos.Group),
		UnAssignGroup: ugr.New(repos.User, repos.Group),
	}
}
