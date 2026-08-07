package annotator

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

var (
	AnnotateImage    = "/ui/annotate/image"
	SubmitBox        = "/ui/annotate/submit-box"
	UpdateBox        = "/ui/annotate/update-box"
	SubmitPolygon    = "/ui/annotate/submit-polygon"
	UpdatePolygon    = "/ui/annotate/update-polygon"
	SubmitImageLabel = "/ui/annotate/submit-label"
	AnnotationPanel  = "/ui/annotate/annotation-panel"
	Annotations      = "/ui/annotate/annotations"
	RemoveAnnotation = "/ui/annotate/remove-annotation"
	SetLabel         = "/ui/annotate/set-label"
)

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(mws...)
		r.Get(AnnotateImage, s.AnnotateImage)

		r.Post(SubmitBox, s.SubmitBox)
		r.Put(UpdateBox, s.UpdateBox)
		r.Post(SubmitPolygon, s.SubmitPolygon)
		r.Put(UpdatePolygon, s.UpdatePolygon)
		r.Post(SubmitImageLabel, s.SubmitLabel)
		r.Get(AnnotationPanel, s.MakeAnnotationPanel)
		r.Get(Annotations, s.GetRegionsAsJSON)
		r.Delete(RemoveAnnotation, s.DeleteAnnotation)
		r.Post(SetLabel, s.SetLabel)
	})
}
