package sqlite

import (
	r "github.com/lejeunel/go-image-annotator/use-cases/role"
	"github.com/lejeunel/go-image-annotator/use-cases/role/create"
	"github.com/lejeunel/go-image-annotator/use-cases/role/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/role/find"
	"github.com/lejeunel/go-image-annotator/use-cases/role/list"
	"github.com/lejeunel/go-image-annotator/use-cases/role/update"
)

func NewSQLiteRoleInteractors(repos *SQLiteRepos) *r.Interactors {
	return &r.Interactors{
		Find:   find.New(repos.Role),
		Create: create.New(repos.Role),
		Delete: delete.New(repos.Role),
		List:   list.New(repos.Role),
		Update: update.New(repos.Role),
	}
}
