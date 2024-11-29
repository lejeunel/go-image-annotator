package tests

import (
	"fmt"
	"github.com/go-test/deep"
	m "go-image-annotator/models"
	"testing"
)

func TestApplyingInvalidBoundingBoxesToImageShoudFail(t *testing.T) {
	tests := map[string]struct {
		xc float64
		yc float64
		w  float64
		h  float64
	}{
		"negative x coord should fail": {xc: -2, yc: 4, h: 3, w: 5},
		"negative y coord should fail": {xc: 2, yc: -4, h: 3, w: 5},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, ctx := NewTestApp(t, 2)

			image := &m.Image{Data: testImage}
			label := &m.Label{Name: "mylabel"}
			collection := &m.Collection{Name: "mylabel"}
			s.Collections.Create(ctx, collection)
			s.Images.Save(ctx, image, collection)
			s.Annotations.CreateLabel(ctx, label)
			bbox := m.NewBoundingBox(tc.xc, tc.yc, tc.h, tc.w)
			bbox.Annotate(label)

			err := s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collection)
			AssertError(t, err)

		})
	}

}

func TestApplyingValidBoundingBoxesToImageShoudSucceed(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: "mylabel"}
	collection := &m.Collection{Name: "mycollection"}
	s.Annotations.CreateLabel(ctx, label)
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	bbox := m.NewBoundingBox(5, 6, 10, 10)
	bbox.Annotate(label)

	err := s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collection)
	AssertNoError(t, err)

	retrievedImage, err := s.Images.Get(ctx, collection.Id.String(), image.Id.String(), true)
	AssertNoError(t, err)

	diff := deep.Equal(image, retrievedImage)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image structs, but got different fields: %v", diff))
	}

}
