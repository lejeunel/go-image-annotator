package image

import (
	"github.com/lejeunel/go-image-annotator-v2/use-cases/image/ingest"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/image/read"
)

type Interactors struct {
	Ingest              ingest.Interactor
	Read                read.Interactor
	List                list.Interactor
	AllowedImageFormats []string
	DefaultPageSize     int
}
