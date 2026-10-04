package annotate

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addply "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
)

type AnnotationPresenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p AnnotationPresenter) SuccessAddBox(r addbox.Response) {
	json.WriteJSON(p.Writer, 200, "successfully added box")
}

func (p AnnotationPresenter) SuccessAddPolygon(r addply.Response) {
	json.WriteJSON(p.Writer, 200, "successfully added polygon")
}

func (p AnnotationPresenter) SuccessDeleteAnnotation(id string) {
	json.WriteJSON(p.Writer, 200, fmt.Sprintf("successfully deleted annotation %v", id))
}

func (p AnnotationPresenter) SuccessUpdateLabel(r updlbl.Response) {
	json.WriteJSON(
		p.Writer,
		200,
		fmt.Sprintf("successfully updated annotation %v with label %v", r.Id, r.Label),
	)
}

func NewAnnotationPresenter(w http.ResponseWriter, l slog.Logger) AnnotationPresenter {
	return AnnotationPresenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
