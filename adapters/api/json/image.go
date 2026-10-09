package json

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ig "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	s "github.com/lejeunel/go-image-annotator/shared"
	"github.com/lejeunel/go-image-annotator/use-cases/image/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/image/raw"
	"github.com/lejeunel/go-image-annotator/use-cases/image/slice"
)

func BuildImageResponse(image im.Image) models.Image {
	response := models.Image{
		Id:         image.Id.String(),
		Collection: image.Collection.Name,
	}
	if len(image.Labels) > 0 {
		labelsToAdd := []string{}
		for _, l := range image.Labels {
			labelsToAdd = append(labelsToAdd, l.Label.Name)
		}
		response.Labels = &labelsToAdd
	}

	if len(image.BoundingBoxes) > 0 {
		boxesToAdd := []models.BoundingBox{}
		for _, b := range image.BoundingBoxes {
			boxesToAdd = append(boxesToAdd,
				models.BoundingBox{
					Id: b.Id.String(),
					Xc: b.Xc, Yc: b.Yc, Height: b.Height, Width: b.Width, Label: b.Label.Name,
				})
		}
		response.BoundingBoxes = &boxesToAdd
	}

	if len(image.Polygons) > 0 {
		polygonsToAdd := []models.Polygon{}
		for _, poly := range image.Polygons {
			points := []models.Point{}
			for _, p := range poly.Points.Coordinates {
				points = append(points, models.Point{p[0], p[1]})
			}
			polygonsToAdd = append(polygonsToAdd,
				models.Polygon{
					Id:     poly.Id.String(),
					Points: points, Label: poly.Label.Name,
				})
		}
		response.Polygons = &polygonsToAdd
	}

	if len(image.Meta) > 0 {
		toAdd := make(map[string]any)
		for _, m := range image.Meta {
			toAdd[m.Key] = m.Value
		}
		response.Meta = &toAdd
	}

	return response
}

type ImagePresenter struct {
	Writer http.ResponseWriter
	ErrorPresenter
}

func (p ImagePresenter) SuccessReadImage(image im.Image) {
	response := BuildImageResponse(image)
	s.WriteJSON(p.Writer, 200, response)
}

func (p ImagePresenter) SuccessDeleteImage(r delete.Response) {
	p.Writer.WriteHeader(http.StatusNoContent)
}

func (p ImagePresenter) SuccessSliceImages(r slice.Response) {
	var images []models.Image
	for _, image := range r.Images {
		images = append(images, BuildImageResponse(image))
	}

	s.WriteJSON(p.Writer, 200, models.ListImages{
		Images:     images,
		Pagination: BuildPaginationResponse(r.Pagination),
	})
}

func (p ImagePresenter) SuccessReadRawImage(r raw.Response) {
	data, err := io.ReadAll(r.Reader)
	if err != nil {
		WriteError(p.Writer, http.StatusInternalServerError, err.Error())
		return
	}
	sum := sha256.Sum256(data)
	etag := `"` + hex.EncodeToString(sum[:]) + `"`
	p.Writer.Header().Set("ETag", etag)
	p.Writer.Header().Set("Content-Type", r.MIMEType)
	p.Writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	p.Writer.Header().Set("Cache-Control", "private, max-age=3600")
	p.Writer.WriteHeader(http.StatusOK)
	p.Writer.Write(data)
}

func (p ImagePresenter) Success(r ig.Response) {
	id := r.ImageId.String()
	response := models.ImageIngestionResponse{
		Id: &id,
	}

	s.WriteJSON(p.Writer, 200, response)
}

func NewImagePresenter(w http.ResponseWriter, l slog.Logger) ImagePresenter {
	return ImagePresenter{Writer: w, ErrorPresenter: NewErrPresenter(w, l)}
}
