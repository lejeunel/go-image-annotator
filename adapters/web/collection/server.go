package collection

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/clone"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/create"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/find"
	listclc "github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/lejeunel/go-image-annotator/use-cases/collection/update"
	listgrp "github.com/lejeunel/go-image-annotator/use-cases/group/list"
	listpr "github.com/lejeunel/go-image-annotator/use-cases/profile/list-all"
)

type Server struct {
	b.PageBuilder
	b.RowURL
	DefaultPageSize   int
	ListCollectionItr listclc.Interactor
	ListGroupItr      listgrp.Interactor
	ListProfileItr    listpr.Interactor
	CreateItr         create.Interactor
	UpdateItr         update.Interactor
	DeleteItr         delete.Interactor
	CloneItr          clone.Interactor
	FindItr           find.Interactor
}

func New(pb b.PageBuilder, defaultPageSize int,
	c create.Interactor, lc listclc.Interactor,
	u update.Interactor,
	d delete.Interactor, cl clone.Interactor, f find.Interactor,
	lg listgrp.Interactor,
	lp listpr.Interactor,
) Server {
	return Server{
		pb, b.NewRowURLWithId(CollectionUrl, resourceUrlFieldName), defaultPageSize,
		lc, lg, lp, c, u, d, cl, f,
	}
}
