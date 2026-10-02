package role

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	r "github.com/lejeunel/go-image-annotator/use-cases/role"
)

type Auth interface {
	ListMethods() []string
}

type Server struct {
	Page   b.PaginatedListBuilder
	RowUrl b.RowURL
	Roles  r.Interactors
	Auth
}

func New(pb b.PageBuilder, rl r.Interactors, auth Auth) Server {
	rolePage := b.NewPaginatedListBuilder(pb, listRolesFields)
	rolePage.ActivateSidebarEntry(PageName)
	return Server{rolePage, b.NewRowURLWithId(RoleRowUrl, resourceUrlFieldName), rl, auth}
}
