package group

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

var GroupRow = "/ui/user"
var CreateUserForm = "/ui/user/new"

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.AdminGroups, s.ListGroups)
		r.Get(GroupRow, s.TableRow)
		r.Delete(GroupRow, s.Delete)
		r.Put(GroupRow, s.Edit)
		r.Get(CreateUserForm, s.CreateForm)
		r.Post(GroupRow, s.Create)
	})
}
