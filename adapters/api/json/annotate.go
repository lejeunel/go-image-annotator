package json

import (
	"log/slog"
	"net/http"

	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addply "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

type AnnotationPresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func (p AnnotationPresenter) SuccessAddBox(r addbox.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p AnnotationPresenter) SuccessAddPolygon(r addply.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p AnnotationPresenter) SuccessDeleteAnnotation(id string) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p AnnotationPresenter) SuccessUpdateLabel(r updlbl.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func NewAnnotationPresenter(w http.ResponseWriter, l slog.Logger) AnnotationPresenter {
	return AnnotationPresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
