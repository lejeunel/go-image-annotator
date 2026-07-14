package collection

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Collections, s.List)
		r.Get(rt.Collection, s.TableRow)
		r.Post(rt.Collection, s.Create)
		r.Delete(rt.Collection, s.Delete)
		r.Put(rt.Collection, s.Edit)
		r.Get(rt.CreateCollectionForm, s.CreateForm)

	})
}
