package sqlite

import (
	"crypto/sha256"

	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	rea "github.com/lejeunel/go-image-annotator/modules/reader"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
)

func NewSQLiteImageInteractors(repos SQLiteRepos, allowedImageFormats []string, pageSize int,
	auth auth.Authorizer) im.Interactors {
	return im.Interactors{
		Ingest: *ingest.New(repos.Image, repos.Collection,
			repos.Label, repos.Annotation,
			repos.FileStore, sha256.New(), rea.ImageSpecsDetector{}, ingest.WithAuth(auth)),
		Find: find.New(repos.ImageStore),
		List: list.New(repos.Image, repos.ImageStore),
	}
}
