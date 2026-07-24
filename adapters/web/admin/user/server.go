package user

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	g "github.com/lejeunel/go-image-annotator/use-cases/group"
	r "github.com/lejeunel/go-image-annotator/use-cases/role"
	u "github.com/lejeunel/go-image-annotator/use-cases/user"
)

type Server struct {
	Page            b.PaginatedListBuilder
	RowUrl          b.RowURL
	Users           u.Interactors
	Roles           r.Interactors
	Groups          g.Interactors
	DefaultPageSize int
}

func New(pb b.PageBuilder, usr u.Interactors, grp g.Interactors, rl r.Interactors, defaultPageSize int) Server {
	userPage := b.NewPaginatedListBuilder(pb, listUsersFields)
	userPage.ActivateSidebarEntry(UserSidebarEntryName)
	return Server{userPage, b.NewRowURL(UserUrl, "id"), usr, rl, grp, defaultPageSize}
}
