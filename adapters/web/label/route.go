package label

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

var Label = "/ui/label"
var CreateLabelForm = "/ui/label/new"

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(rt.Labels, s.List)
		r.Get(Label, s.TableRow)
		r.Post(Label, s.Create)
		r.Delete(Label, s.Delete)
		r.Put(Label, s.Edit)
		r.Get(CreateLabelForm, s.CreateForm)

	})
}
