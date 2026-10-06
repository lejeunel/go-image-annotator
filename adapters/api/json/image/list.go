package image

import (
	"log/slog"
	"net/http"

	"github.com/lejeunel/go-image-annotator/adapters/api/json"
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/image/slice"
)

type List struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p List) SuccessSliceImages(r slice.Response) {
	var images []models.Image
	for _, image := range r.Images {
		images = append(images, BuildImageResponse(image))
	}

	s.WriteJSON(p.Writer, 200, models.ListImages{
		Images:     images,
		Pagination: json.BuildPaginationResponse(r.Pagination),
	})
}

func NewListPresenter(w http.ResponseWriter, l slog.Logger) List {
	return List{Writer: w, ErrorPresenter: json.NewErrPresenter(w, l)}
}
