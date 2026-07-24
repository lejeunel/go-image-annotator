package image

import (
	"github.com/go-chi/chi/v5"
	rt "github.com/lejeunel/go-image-annotator/routes"
	"net/http"
)

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler) {

	r.Group(func(r chi.Router) {
		r.Use(mws...)

		r.Get(rt.ImagesUrl, s.List)
		r.Get(ImageRow, s.TableRow)
		r.Delete(ImageRow, s.Delete)
	})
}
