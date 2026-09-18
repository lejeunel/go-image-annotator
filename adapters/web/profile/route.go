package profile

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
)

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.ProfilesUrl, s.List)
		r.Get(ProfileUrl, s.TableRow)
		r.Post(ProfileUrl, s.Create)
		r.Delete(ProfileUrl, s.Delete)
		r.Put(ProfileUrl, s.Edit)
		r.Get(CreateProfileFormUrl, s.CreateForm)
	})
}
