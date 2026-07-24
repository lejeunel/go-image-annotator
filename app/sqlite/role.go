package sqlite

import (
	ri "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/role"
	r "github.com/lejeunel/go-image-annotator/use-cases/role"
	"github.com/lejeunel/go-image-annotator/use-cases/role/create"
	"github.com/lejeunel/go-image-annotator/use-cases/role/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/role/find"
	"github.com/lejeunel/go-image-annotator/use-cases/role/list"
	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
)

func NewSQLiteRoleInteractors(repo ri.SQLiteRoleRepo) r.Interactors {
	return r.Interactors{
		Find:   find.New(repo),
		Create: create.New(repo),
		Delete: delete.New(repo),
		List:   list.New(repo),
		Update: update.New(repo),
	}
}
