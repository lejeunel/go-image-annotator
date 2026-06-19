package sqlite

import (
	a "github.com/lejeunel/go-image-annotator/modules/authentifier"
	grp "github.com/lejeunel/go-image-annotator/use-cases/group"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/read"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func NewSQLiteGroupInteractors(repos *SQLiteRepos, tokenGen a.AuthGenerator) *grp.Interactors {
	return &grp.Interactors{
		Find:   read.New(repos.Group),
		Create: create.New(repos.Group),
		Delete: delete.New(repos.Group),
		List:   list.New(repos.Group),
		Update: update.New(repos.Group),
	}
}
