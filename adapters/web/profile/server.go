package profile

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	listlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/create"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/find"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/list"
	"github.com/lejeunel/go-image-annotator/use-cases/profile/update"
)

type Server struct {
	b.PageBuilder
	b.RowURL
	DefaultPageSize  int
	ListItr          list.Interactor
	ListAllLabelsItr listlbl.Interactor
	CreateItr        create.Interactor
	UpdateItr        update.Interactor
	DeleteItr        delete.Interactor
	FindItr          find.Interactor
}

func New(
	pb b.PageBuilder,
	defaultPageSize int,
	c create.Interactor,
	l list.Interactor,
	llbl listlbl.Interactor,
	u update.Interactor,
	d delete.Interactor,
	f find.Interactor,
) Server {
	return Server{
		pb,
		b.NewRowURLWithId(ProfileUrl, resourceUrlFieldName),
		defaultPageSize,
		l,
		llbl,
		c,
		u,
		d,
		f,
	}
}
