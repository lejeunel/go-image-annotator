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
		r.Get(LabelUrl, s.TableRow)
		r.Post(LabelUrl, s.Create)
		r.Delete(LabelUrl, s.Delete)
		r.Put(LabelUrl, s.Edit)
		r.Get(CreateLabelFormUrl, s.CreateForm)

	})
}
