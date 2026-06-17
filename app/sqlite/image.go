package sqlite

import (
	"crypto/sha256"
	rea "github.com/lejeunel/go-image-annotator/modules/reader"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/read"
)

func NewSQLiteImageInteractors(repos SQLiteRepos, allowedImageFormats []string, pageSize int) im.Interactors {
	return im.Interactors{
		Ingest: *ingest.New(repos.Image, repos.Collection,
			repos.Label, repos.Annotation,
			repos.FileStore, sha256.New(), rea.ImageSpecsDetector{}),
		Read:                read.New(repos.ImageStore),
		List:                list.New(repos.Image, repos.ImageStore),
		AllowedImageFormats: allowedImageFormats,
		DefaultPageSize:     pageSize,
	}
}
