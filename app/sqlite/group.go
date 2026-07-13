package sqlite

import (
	gi "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	grp "github.com/lejeunel/go-image-annotator/use-cases/group"
	"github.com/lejeunel/go-image-annotator/use-cases/group/create"
	"github.com/lejeunel/go-image-annotator/use-cases/group/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/group/find"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

func NewSQLiteGroupInteractors(repo gi.SQLiteGroupRepo) grp.Interactors {
	return grp.Interactors{
		Find:   find.New(repo),
		Create: create.New(repo),
		Delete: delete.New(repo),
		List:   list.New(repo),
		Update: update.New(repo),
	}
}
