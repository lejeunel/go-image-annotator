package sqlite

import (
	"github.com/jonboulle/clockwork"
	ar "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	cr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	gr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/group"
	ir "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	q "github.com/lejeunel/go-image-annotator/modules/job-queue"
	"github.com/lejeunel/go-image-annotator/shared/validation"
	clc "github.com/lejeunel/go-image-annotator/use-cases/collection"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/find"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
	"log/slog"
)

func NewSQLiteCollectionInteractors(
	cr cr.SQLiteCollectionRepo,
	ir ir.SQLiteImageRepo,
	ar ar.SQLiteAnnotationRepo,
	gr gr.SQLiteGroupRepo,
	ims ims.ImageStore,
	el el.EventLogger,
	logger slog.Logger,
	pageSize int, auth auth.Interface) clc.Interactors {
	return clc.Interactors{
		Find: find.New(cr),
		Create: create.New(cr, gr, create.WithNameValidator(validation.NewNameValidator()),
			create.WithClock(clockwork.NewRealClock()), create.WithAuth(auth)),
		Delete: delete.New(cr, gr, delete.WithAuth(auth)),
		List:   list.New(cr),
		Update: update.New(cr, gr, update.WithAuth(auth)),
		Clone:  clone.New(ir, cr, ar, gr, ims, el, logger, q.NewJobQueue()),
	}
}
