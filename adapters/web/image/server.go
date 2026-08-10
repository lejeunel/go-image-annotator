package image

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	ia "github.com/lejeunel/go-image-annotator/use-cases/image/ingest-archive"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
)

type Server struct {
	b.PageBuilder
	DefaultPageSize  int
	ListItr          list.Interactor
	DeleteItr        delete.Interactor
	FindItr          find.Interactor
	IngestArchiveItr ia.Interactor
}

func New(
	pb b.PageBuilder, defaultPageSize int,
	l list.Interactor, d delete.Interactor, f find.Interactor,
	i ia.Interactor,
) Server {
	return Server{pb, defaultPageSize, l, d, f, i}
}
