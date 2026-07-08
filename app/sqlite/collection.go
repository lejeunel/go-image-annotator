package sqlite

import (
	"github.com/jonboulle/clockwork"
	ci "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	gi "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/shared/validation"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/find"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func NewSQLiteCollectionInteractors(cr ci.SQLiteCollectionRepo,
	gr gi.SQLiteGroupRepo,
	pageSize int, auth auth.Authorizer) clc.Interactors {
	return clc.Interactors{
		Find: find.New(cr),
		Create: create.New(cr, gr, create.WithNameValidator(validation.NewNameValidator()),
			create.WithClock(clockwork.NewRealClock()), create.WithAuth(auth)),
		Delete: delete.New(cr, gr, delete.WithAuth(auth)),
		List:   list.New(cr),
		Update: update.New(cr, gr, update.WithAuth(auth)),
	}
}
