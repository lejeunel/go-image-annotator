package annotator

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var AnnotateImage = "/ui/annotate/image"
var SubmitBox = "/ui/annotate/submit-box"
var UpdateBox = "/ui/annotate/update-box"
var SubmitPolygon = "/ui/annotate/submit-polygon"
var UpdatePolygon = "/ui/annotate/update-polygon"
var SubmitImageLabel = "/ui/annotate/submit-label"
var AnnotationPanel = "/ui/annotate/annotation-panel"
var Annotations = "/ui/annotate/annotations"
var RemoveAnnotation = "/ui/annotate/remove-annotation"
var SetLabel = "/ui/annotate/set-label"

func (s *Server) Route(r chi.Router,
	mws ...func(http.Handler) http.Handler) {

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
