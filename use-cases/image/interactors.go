package image

import (
	"github.com/lejeunel/go-image-annotator/modules/auth"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/read"
)

type Interactors struct {
	Ingest              ingest.Interactor
	Read                read.Interactor
	List                list.Interactor
	AllowedImageFormats []string
	DefaultPageSize     int
	Authorizer          auth.Auth
}
