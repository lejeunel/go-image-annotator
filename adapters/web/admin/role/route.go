package role

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.AdminRolesUrl, s.ListRoles)
		r.Get(RoleRowUrl, s.TableRow)
		r.Delete(RoleRowUrl, s.Delete)
		r.Put(RoleRowUrl, s.Edit)
		r.Get(CreateRoleForm, s.CreateForm)
		r.Post(RoleRowUrl, s.Create)
	})
}
