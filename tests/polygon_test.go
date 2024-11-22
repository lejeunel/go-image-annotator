package tests

import (
	"fmt"
	"github.com/go-test/deep"
	m "go-image-annotator/models"
	"testing"
)

func TestApplyingPolygonToImage(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image, _ := s.Images.Save(ctx, &m.Image{Data: testImage})
	label, _ := s.Annotations.Create(ctx, &m.Label{Name: "mylabel"})
	polygon, err := m.NewBoundingBox(10, 10, 30, 30)
	AssertNoError(t, err)
	polygon.Label = label

	image, err = s.Annotations.ApplyPolygonToImage(ctx, polygon, image)
	AssertNoError(t, err)

	retrievedImage, err := s.Images.GetOne(ctx, image.Id.String(), true)
	AssertNoError(t, err)

	diff := deep.Equal(image, retrievedImage)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image structs, but got different fields: %v", diff))
	}

}

func TestDeletePolygon(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image := &m.Image{Data: testImage}
	image, err := s.Images.Save(ctx, image)

	polygon, err := m.NewBoundingBox(10, 10, 30, 30)

	image, err = s.Annotations.ApplyPolygonToImage(ctx, polygon, image)
	image, err = s.Annotations.DeletePolygonFromImage(ctx, polygon, image)
	AssertNoError(t, err)
	image, _ = s.Images.GetOne(ctx, image.Id.String(), false)

	nPolygons := len(image.Polygons)
	if nPolygons != 0 {
		t.Fatalf("expected to retrieve image without polygons, but got %v.", nPolygons)
	}

}
