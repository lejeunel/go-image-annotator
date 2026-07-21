package image

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/find"
	"github.com/lejeunel/go-image-annotator/use-cases/image/list"
)

type Server struct {
	b.PageBuilder
	DefaultPageSize int
	ListItr         list.Interactor
	DeleteItr       delete.Interactor
	FindItr         find.Interactor
}

func New(b b.PageBuilder, defaultPageSize int,
	l list.Interactor, d delete.Interactor, f find.Interactor) Server {
	return Server{b, defaultPageSize, l, d, f}
}
