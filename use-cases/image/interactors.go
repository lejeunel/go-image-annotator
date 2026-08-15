package image

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/ingest"
	aig "github.com/lejeunel/go-image-annotator/use-cases/image/ingest-archive"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
	"github.com/lejeunel/go-image-annotator/use-cases/image/scroll"
)

type Interactors struct {
	Ingest          ingest.Interactor
	IngestArchive   aig.Interactor
	Find            find.Interactor
	List            list.Interactor
	Scroll          scroll.Interactor
	Raw             raw.Interactor
	Delete          delete.Interactor
	DefaultPageSize int
	Authorizer      auth.Authorizer
}
