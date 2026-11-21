package server

import (
	"fmt"

	"github.com/lejeunel/go-image-annotator/adapters/api/models"
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
)

// pointsFromModel converts API points into domain points. Each point must
// hold exactly two coordinates (x, y): the generated models do not enforce
// the spec's minItems/maxItems, so a malformed point would otherwise cause
// an out-of-range index.
func pointsFromModel(points []models.Point) (an.Points, error) {
	var res an.Points
	for i, point := range points {
		if len(point) != 2 {
			return an.Points{}, fmt.Errorf(
				"point %v must have exactly 2 coordinates, got %v", i, len(point))
		}
		res.Append(point[0], point[1])
	}
	return res, nil
}
