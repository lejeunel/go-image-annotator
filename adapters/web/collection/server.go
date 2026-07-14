package collection

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/find"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
)

type Server struct {
	b.PageBuilder
	DefaultPageSize int
	ListItr         list.Interactor
	CreateItr       create.Interactor
	UpdateItr       update.Interactor
	DeleteItr       delete.Interactor
	FindItr         find.Interactor
}

func New(b b.PageBuilder, defaultPageSize int,
	c create.Interactor, l list.Interactor, u update.Interactor, d delete.Interactor, f find.Interactor) Server {
	return Server{b, defaultPageSize, l, c, u, d, f}
}
