package image

import (
	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	im "github.com/lejeunel/go-image-annotator/entities/image"
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
