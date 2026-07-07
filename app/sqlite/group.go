package sqlite

import (
	grp "github.com/lejeunel/go-image-annotator/use-cases/group"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/group/find"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func NewSQLiteGroupInteractors(repos *SQLiteRepos) *grp.Interactors {
	return &grp.Interactors{
		Find:   find.New(repos.Group),
		Create: create.New(repos.Group),
		Delete: delete.New(repos.Group),
		List:   list.New(repos.Group),
		Update: update.New(repos.Group),
	}
}
