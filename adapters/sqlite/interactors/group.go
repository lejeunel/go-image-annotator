package interactors

import (
	infra "github.com/lejeunel/go-image-annotator/adapters/sqlite"
	tok "github.com/lejeunel/go-image-annotator/app/token"
	grp "github.com/lejeunel/go-image-annotator/use-cases/group"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/read"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func NewSQLiteGroupInteractors(repos *infra.SQLiteInfra, tokenGen tok.TokenGenerator) *grp.Interactors {
	return &grp.Interactors{
		Find:   read.New(repos.Group),
		Create: create.New(repos.Group),
		Delete: delete.New(repos.Group),
		List:   list.New(repos.Group),
		Update: update.New(repos.Group),
	}
}
