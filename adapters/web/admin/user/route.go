package user

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Admin, s.ListUsers)
		r.Get(rt.AdminUsers, s.ListUsers)
		r.Get(UserUrl, s.TableRow)
		r.Delete(UserUrl, s.Delete)
		r.Get(CreateUserFormUrl, s.CreateForm)
		r.Post(UserUrl, s.Create)
		r.Put(UserUrl, s.Edit)
	})
}
