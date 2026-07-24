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
		r.Get(CollectionUrl, s.TableRow)
		r.Post(CollectionUrl, s.Create)
		r.Delete(CollectionUrl, s.Delete)
		r.Put(CollectionUrl, s.Edit)
		r.Get(CreateCollectionFormUrl, s.CreateForm)

	})
}
