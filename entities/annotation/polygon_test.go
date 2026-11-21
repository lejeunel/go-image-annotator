package annotation

import (
	"math"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func pts(coords ...[2]float32) Points {
	return Points{Coordinates: coords}
}

func TestPolygonArea(t *testing.T) {
	tests := []struct {
		name   string
		points Points
		want   float64
	}{
		{
			"unit square",
			pts([2]float32{0, 0}, [2]float32{1, 0}, [2]float32{1, 1}, [2]float32{0, 1}),
			1,
		},
		{
			"unit square, clockwise",
			pts([2]float32{0, 0}, [2]float32{0, 1}, [2]float32{1, 1}, [2]float32{1, 0}),
			1,
		},
		{"right triangle", pts([2]float32{0, 0}, [2]float32{4, 0}, [2]float32{0, 3}), 6},
		{
			"concave L-shape",
			pts([2]float32{0, 0}, [2]float32{2, 0}, [2]float32{2, 1},
				[2]float32{1, 1}, [2]float32{1, 2}, [2]float32{0, 2}),
			3,
		},
		{"collinear", pts([2]float32{0, 0}, [2]float32{1, 1}, [2]float32{2, 2}), 0},
		{
			"symmetric bow-tie",
			pts([2]float32{0, 0}, [2]float32{1, 1}, [2]float32{1, 0}, [2]float32{0, 1}),
			0,
		},
		{"empty", Points{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, tt.points.Area(), 1e-9)
		})
	}
}

func TestValidatePolygon(t *testing.T) {
	nan := float32(math.NaN())
	inf := float32(math.Inf(1))
	tests := []struct {
		name    string
		points  Points
		wantErr bool
	}{
		{"no points", Points{}, true},
		{"one point", pts([2]float32{0, 0}), true},
		{"two points", pts([2]float32{0, 0}, [2]float32{1, 1}), true},
		{"collinear points", pts([2]float32{0, 0}, [2]float32{1, 1}, [2]float32{2, 2}), true},
		{
			"collinear points with non-integer coordinates",
			pts([2]float32{0.1, 0.1}, [2]float32{0.2, 0.2}, [2]float32{0.3, 0.3}),
			true,
		},
		{"repeated point", pts([2]float32{1, 1}, [2]float32{1, 1}, [2]float32{1, 1}), true},
		{"NaN coordinate", pts([2]float32{0, 0}, [2]float32{nan, 0}, [2]float32{0, 1}), true},
		{"infinite coordinate", pts([2]float32{0, 0}, [2]float32{inf, 0}, [2]float32{0, 1}), true},
		{"triangle", pts([2]float32{0, 0}, [2]float32{1, 0}, [2]float32{0, 1}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.points)
			if tt.wantErr {
				assert.ErrorIs(t, err, e.ErrValidation)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
