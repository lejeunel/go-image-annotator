package group

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	g "github.com/lejeunel/go-image-annotator/use-cases/group"
)

type Server struct {
	Page   b.PaginatedListBuilder
	RowUrl b.RowURL
	Groups g.Interactors
}

func New(pb b.PageBuilder, grp g.Interactors) Server {
	groupPage := b.NewPaginatedListBuilder(pb, listGroupsFields)
	groupPage.ActivateSidebarEntry(GroupPageEntryName)
	return Server{groupPage, b.NewRowURL(GroupRowUrl, resourceUrlFieldName), grp}
}
