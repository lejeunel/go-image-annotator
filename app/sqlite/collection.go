package sqlite

import (
	"github.com/jonboulle/clockwork"
	ci "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	gi "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	"github.com/lejeunel/go-image-annotator/shared/validation"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/read"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func NewSQLiteCollectionInteractors(cr ci.SQLiteCollectionRepo,
	gr gi.SQLiteGroupRepo,
	pageSize int) clc.Interactors {
	return clc.Interactors{
		Find: read.New(cr),
		Create: create.New(cr, gr, create.WithNameValidator(validation.NewNameValidator()),
			create.WithClock(clockwork.NewRealClock())),
		Delete:          delete.New(cr, gr),
		List:            list.New(cr),
		Update:          update.New(cr, gr),
		DefaultPageSize: pageSize,
	}
}
