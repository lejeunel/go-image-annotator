package image

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
)

type Interactors struct {
	Ingest              ingest.Interactor
	Find                find.Interactor
	List                list.Interactor
	Raw                 raw.Interactor
	AllowedImageFormats []string
	DefaultPageSize     int
	Authorizer          auth.Authorizer
}
