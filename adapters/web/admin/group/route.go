package group

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router, mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.AdminGroupsUrl, s.ListGroups)
		r.Get(GroupRowUrl, s.TableRow)
		r.Delete(GroupRowUrl, s.Delete)
		r.Put(GroupRowUrl, s.Edit)
		r.Get(CreateUserFormUrl, s.CreateForm)
		r.Post(GroupRowUrl, s.Create)
	})
}
