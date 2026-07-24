package role

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	r "github.com/lejeunel/go-image-annotator/use-cases/role"
)

type Server struct {
	Page   b.PaginatedListBuilder
	RowUrl b.RowURL
	Roles  r.Interactors
}

func New(pb b.PageBuilder, rl r.Interactors) Server {
	rolePage := b.NewPaginatedListBuilder(pb, listRolesFields)
	rolePage.ActivateSidebarEntry(RolePageEntryName)
	return Server{rolePage, b.NewRowURL(RoleRowUrl, resourceUrlFieldName), rl}
}
