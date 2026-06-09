package interactors

import (
	"github.com/jonboulle/clockwork"
	ci "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/collection"
	gi "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/group"
	"github.com/lejeunel/go-image-annotator/shared/validation"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/read"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func NewSQLiteCollectionInteractors(cr *ci.SQLiteCollectionRepo,
	gr *gi.SQLiteGroupRepo,
	pageSize int) *clc.Interactors {
	return &clc.Interactors{
		Find: read.NewInteractor(cr),
		Create: *create.NewInteractor(cr, gr, create.WithNameValidator(validation.NewNameValidator()),
			create.WithClock(clockwork.NewRealClock())),
		Delete:          delete.NewInteractor(cr, gr),
		List:            list.NewInteractor(cr),
		Update:          update.NewInteractor(cr, gr),
		DefaultPageSize: pageSize,
	}
}
