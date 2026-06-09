package interactors

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/sqlite"
	tg "github.com/lejeunel/go-image-annotator/app/token-generator"
	grp "github.com/lejeunel/go-image-annotator/use-cases/group"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/read"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func NewSQLiteGroupInteractors(repos *infra.SQLiteInfra, tokenGen tg.TokenGenerator) *grp.Interactors {
	return &grp.Interactors{
		Find:   read.NewInteractor(repos.GroupRepo),
		Create: create.NewInteractor(repos.GroupRepo),
		Delete: delete.NewInteractor(repos.GroupRepo),
		List:   list.NewInteractor(repos.GroupRepo),
		Update: update.NewInteractor(repos.GroupRepo),
	}
}
