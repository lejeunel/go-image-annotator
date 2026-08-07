package sqlite

import (
	anrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/annotation"
	clrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	imrepo "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/image"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	ingm "github.com/lejeunel/go-image-annotator/modules/ingester"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
)

func NewSQLiteImageInteractors(
	imr imrepo.SQLiteImageRepo,
	clr clrepo.SQLiteCollectionRepo,
	anr anrepo.SQLiteAnnotationRepo,
	ims ims.ImageStore,
	imfs fs.Interface,
	ingester ingm.Interface,
	pageSize int,
	auth auth.Interface,
) im.Interactors {
	return im.Interactors{
		Ingest: *ingest.New(ingester, clr, ingest.WithAuth(auth)),
		Find:   find.New(ims),
		Raw:    raw.New(imfs, imr),
		List:   list.New(imr, ims),
		Delete: delete.New(ims, imr, anr),
	}
}
