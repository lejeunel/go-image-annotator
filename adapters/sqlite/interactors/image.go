package interactors

import (
	"crypto/sha256"
	"github.com/lejeunel/go-image-annotator/adapters/infra"
	rea "github.com/lejeunel/go-image-annotator/app/reader"
	im "github.com/lejeunel/go-image-annotator/use-cases/image"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/read"
)

func NewSQLiteImageInteractors(repos *infra.SQLiteInfra, allowedImageFormats []string) *im.Interactors {
	return &im.Interactors{
		Ingest: *ingest.NewInteractor(repos.ImageRepo, repos.CollectionRepo,
			repos.LabelRepo, repos.AnnotationRepo,
			repos.FileStore, sha256.New(), rea.ImageSpecsDetector{}),
		Read:                *read.NewInteractor(*repos.ImageStore),
		List:                *list.NewInteractor(repos.ImageRepo, repos.ImageStore),
		AllowedImageFormats: allowedImageFormats,
		DefaultPageSize:     10,
	}
}
