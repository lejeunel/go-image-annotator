package image

import (
	"github.com/lejeunel/go-image-annotator-v2/adapters/api/json"
	"github.com/lejeunel/go-image-annotator-v2/adapters/api/models"
	"github.com/lejeunel/go-image-annotator-v2/use-cases/image/list"
	"net/http"
)

type List struct {
	Writer http.ResponseWriter
	json.ErrorPresenter
}

func (p List) Success(r list.Response) {
	response := models.ListImagesResponse{
		Pagination: json.BuildPaginationResponse(r.Pagination),
	}

	for _, image := range r.Images {
		response.Images = append(response.Images, BuildImageResponse(image))
	}

	json.WriteJSON(p.Writer, 200, response)

}

func NewListPresenter(w http.ResponseWriter) List {
	return List{Writer: w, ErrorPresenter: json.ErrorPresenter{Writer: w}}
}
