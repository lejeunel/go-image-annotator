package user

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

var User = "/ui/user"
var CreateUserForm = "/ui/user/new"

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Admin, s.ListUsers)
		r.Get(rt.AdminUsers, s.ListUsers)
		r.Get(User, s.TableRow)
		r.Delete(User, s.Delete)
		r.Get(CreateUserForm, s.CreateForm)
		r.Post(User, s.Create)
	})
}
