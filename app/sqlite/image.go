package sqlite

import (
	anrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	imrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	ing "github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
)

func NewSQLiteImageInteractors(
	imr imrepo.SQLiteImageRepo,
	clr clrepo.SQLiteCollectionRepo,
	anr anrepo.SQLiteAnnotationRepo,
	ims ims.ImageStore,
	imfs fs.FileStore,
	ingester ing.Ingester,
	pageSize int,
	auth auth.Interface,
) im.Interactors {
	return im.Interactors{
		Ingest: *ing.New(ingester, clr, ing.WithAuth(auth)),
		Find:   find.New(ims),
		Raw:    raw.New(imfs, imr),
		List:   list.New(imr, ims),
		Delete: delete.New(ims),
	}
}
