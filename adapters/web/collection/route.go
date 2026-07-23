package collection

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

var Collection = "/ui/collection"
var CreateCollectionForm = "/ui/collection/new"

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Collections, s.List)
		r.Get(Collection, s.TableRow)
		r.Post(Collection, s.Create)
		r.Delete(Collection, s.Delete)
		r.Put(Collection, s.Edit)
		r.Get(CreateCollectionForm, s.CreateForm)

	})
}
