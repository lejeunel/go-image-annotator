package user

import (
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	u "github.com/lejeunel/go-image-annotator/use-cases/user"
)

type Server struct {
	Page            b.PaginatedListBuilder
	Users           u.Interactors
	DefaultPageSize int
}

func New(pb b.PageBuilder, usr u.Interactors, defaultPageSize int) Server {
	userPage := b.NewPaginatedListBuilder(pb, listUsersFields)
	return Server{userPage, usr, defaultPageSize}
}
