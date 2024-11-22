package tests

import (
	"fmt"
	"github.com/go-test/deep"
	m "go-image-annotator/models"
	"testing"
)

func TestCreatingInvalidLabelShouldFail(t *testing.T) {
	tests := map[string]struct {
		name string
	}{
		"with spaces":         {name: "the name with spaces"},
		"with capitals":       {name: "LaBeL NaMe"},
		"with specials chars": {name: "l4b3l n4m3"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, ctx := NewTestApp(t, 2)
			label := &m.Label{Name: tc.name}
			label, err := s.Annotations.Create(ctx, label)
			AssertError(t, err)
		})
	}
}

func TestCreateAndRetrieveLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel",
		Description: "the description"}

	label, err := s.Annotations.Create(ctx, label)
	AssertNoError(t, err)

	retrievedLabel, err := s.Annotations.GetOne(ctx, label.Id.String())

	diff := deep.Equal(label, retrievedLabel)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image structs, but got different fields: %v", diff))
	}

}

func TestDeleteLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}

	label, err := s.Annotations.Create(ctx, label)
	err = s.Annotations.Delete(ctx, label)

	AssertNoError(t, err)

	label, err = s.Annotations.GetOne(ctx, label.Id.String())
	AssertError(t, err)

}

func TestDeletingUsedLabelShouldFail(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}

	label, err := s.Annotations.Create(ctx, label)
	AssertNoError(t, err)
	image := &m.Image{Data: testImage}
	image, err = s.Annotations.ApplyLabelToImage(ctx, label, image)
	AssertNoError(t, err)

	err = s.Annotations.Delete(ctx, label)
	AssertError(t, err)

}

func TestDeleteLabeledImageAndItsAssociatedLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}
	label, _ = s.Annotations.Create(ctx, label)
	image := &m.Image{Data: testImage}
	image, _ = s.Annotations.ApplyLabelToImage(ctx, label, image)

	s.Annotations.Delete(ctx, label)
	s.Images.Delete(ctx, image)
	numAssociatedImages, _ := s.Images.LabelRepo.NumImagesWithLabel(ctx, label)

	if numAssociatedImages != 0 {
		t.Fatalf("expected to have 0 associated images, got %v", numAssociatedImages)
	}

}

func TestApplyingLabelToImage(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image, err := s.Images.Save(ctx, &m.Image{Data: testImage})
	label, err := s.Annotations.Create(ctx, &m.Label{Name: "mylabel"})
	image, err = s.Annotations.ApplyLabelToImage(ctx, label, image)
	AssertNoError(t, err)

	retrievedImage, err := s.Images.GetOne(ctx, image.Id.String(), false)
	AssertNoError(t, err)
	nLabels := len(retrievedImage.Annotations)
	if len(retrievedImage.Annotations) != 1 {
		t.Fatalf("expected to retrieve image with 1 label, but got %v.", nLabels)
	}

}

func TestRemovingLabelFromImage(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image, err := s.Images.Save(ctx, &m.Image{Data: testImage})
	label, err := s.Annotations.Create(ctx, &m.Label{Name: "mylabel"})
	image, err = s.Annotations.ApplyLabelToImage(ctx, label, image)
	image, err = s.Annotations.RemoveAnnotationFromImage(ctx, image.Annotations[0], image)
	AssertNoError(t, err)

	if len(image.Annotations) != 0 {
		t.Fatalf("expected to retrieve image with 0 label, but got %v", len(image.Annotations))
	}

}

func TestApplyingPolygonToImage(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: "mylabel"}

	image, _ = s.Images.Save(ctx, image)
	label, _ = s.Annotations.Create(ctx, label)

	polygon, err := m.NewBoundingBox(10, 10, 30, 30)
	AssertNoError(t, err)
	polygon.Label = label

	image, err = s.Annotations.ApplyPolygonToImage(ctx, polygon, image)
	AssertNoError(t, err)

	retrievedImage, err := s.Images.GetOne(ctx, image.Id.String(), false)
	AssertNoError(t, err)

	polygons := retrievedImage.Polygons
	if len(polygons) != 1 {
		t.Fatalf("expected to retrieve image with 1 polygon, but got %v.", len(polygons))
	}

	polygonLabel := polygons[0].Label.Name
	if polygonLabel != "mylabel" {
		t.Fatalf("expected to retrieve polygon with label mylabel, but got %v.", polygonLabel)
	}

}

func TestInvalidBoundingBoxesShouldFail(t *testing.T) {
	tests := map[string]struct {
		x0 int
		y0 int
		x1 int
		y1 int
	}{
		"negative values":   {x0: -2, y0: 4, x1: 5, y1: 9},
		"inverted x values": {x0: 5, y0: 2, x1: 0, y1: 4},
		"inverted y values": {x0: 0, y0: 9, x1: 5, y1: 4},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, ctx := NewTestApp(t, 2)
			image := &m.Image{Data: testImage}

			image, _ = s.Images.Save(ctx, image)
			_, err := m.NewBoundingBox(tc.x0, tc.y0, tc.x1, tc.y1)

			AssertError(t, err)
		})
	}
}

func TestDeletePolygon(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image := &m.Image{Data: testImage}
	image, err := s.Images.Save(ctx, image)

	polygon, err := m.NewBoundingBox(10, 10, 30, 30)

	image, err = s.Annotations.ApplyPolygonToImage(ctx, polygon, image)

	err = s.Annotations.DeletePolygon(ctx, polygon)
	AssertNoError(t, err)
	image, _ = s.Images.GetOne(ctx, image.Id.String(), false)

	nPolygons := len(image.Polygons)
	if nPolygons != 0 {
		t.Fatalf("expected to retrieve image without polygons, but got %v.", nPolygons)
	}

}
