package interactors

import (
	"github.com/jonboulle/clockwork"
	infra "github.com/lejeunel/go-image-annotator/infra/db/collection"
	"github.com/lejeunel/go-image-annotator/shared/validation"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/read"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

func NewSQLiteCollectionInteractors(repo *infra.SQLiteCollectionRepo, pageSize int) *clc.Interactors {
	return &clc.Interactors{
		Find: read.NewInteractor(repo),
		Create: *create.NewInteractor(repo, create.WithNameValidator(validation.NewNameValidator()),
			create.WithClock(clockwork.NewRealClock())),
		Delete:          delete.NewInteractor(repo),
		List:            list.NewInteractor(repo),
		Update:          update.NewInteractor(repo),
		DefaultPageSize: pageSize,
	}
}
