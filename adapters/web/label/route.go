package label

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Labels, s.List)
		r.Get(rt.Label, s.TableRow)
		r.Post(rt.Label, s.Create)
		r.Delete(rt.Label, s.Delete)
		r.Put(rt.Label, s.Edit)
		r.Get(rt.CreateLabelForm, s.CreateForm)

	})
}
