package sqlite

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ingm "github.com/lejeunel/go-image-annotator/modules/ingester"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
)

func NewSQLiteImageInteractors(repos SQLiteRepos, ingester ingm.Interface, pageSize int,
	auth auth.Authorizer) im.Interactors {
	return im.Interactors{
		Ingest: *ingest.New(ingester, repos.Collection, ingest.WithAuth(auth)),
		Find:   find.New(repos.ImageStore),
		Raw:    raw.New(repos.FileStore, repos.Image),
		List:   list.New(repos.Image, repos.ImageStore),
		Delete: delete.New(repos.ImageStore, repos.Image, repos.Annotation),
	}
}
