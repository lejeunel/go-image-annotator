package annotate

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
)

type AnnotationPresenter struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p AnnotationPresenter) SuccessAddBox(r addbox.Response) {
	json.WriteJSON(p.Writer, 200, "successfully added box")
}

func NewAnnotationPresenter(w http.ResponseWriter, l slog.Logger) AnnotationPresenter {
	return AnnotationPresenter{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
