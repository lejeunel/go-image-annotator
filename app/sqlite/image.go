package sqlite

import (
	"log/slog"

	anrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	imrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	q "github.com/lejeunel/go-image-annotator/modules/job-queue"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	ing "github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	ia "github.com/lejeunel/go-image-annotator/use-cases/image/ingest-archive"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
)

func NewSQLiteImageInteractors(
	imr imrepo.SQLiteImageRepo,
	clr clrepo.SQLiteCollectionRepo,
	anr anrepo.SQLiteAnnotationRepo,
	ims ims.ImageStore,
	imfs fs.FileStore,
	tmpfs fs.LocalFileStore,
	imageIngester ing.Ingester,
	archiveIngester ia.ArchiveIngester,
	el el.EventLogger,
	logger slog.Logger,
	pageSize int,
	auth auth.Interface,
) im.Interactors {
	return im.Interactors{
		Ingest:        *ing.New(imageIngester, clr, ing.WithAuth(auth)),
		IngestArchive: ia.New(archiveIngester, clr, tmpfs, el, logger, q.NewAsyncJobQueue(), ia.WithAuth(auth)),
		Find:          find.New(ims),
		Raw:           raw.New(imfs, imr),
		List:          list.New(imr, ims),
		Delete:        delete.New(ims),
	}
}
