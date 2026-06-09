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
		Find:          *read.NewInteractor(repos.UserRepo),
		Create:        *create.NewInteractor(repos.UserRepo, tokenGen),
		Delete:        *delete.NewInteractor(repos.UserRepo),
		List:          *list.NewInteractor(repos.UserRepo),
		RenewToken:    *rt.NewInteractor(repos.UserRepo, tokenGen),
		AssignRole:    *aro.NewInteractor(repos.UserRepo),
		UnAssignRole:  *uar.NewInteractor(repos.UserRepo),
		AssignGroup:   *agr.NewInteractor(repos.UserRepo, repos.GroupRepo),
		UnAssignGroup: *ugr.NewInteractor(repos.UserRepo, repos.GroupRepo),
	}
}
