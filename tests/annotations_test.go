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
			err := s.Annotations.CreateLabel(ctx, label)
			AssertError(t, err)
		})
	}
}

func TestCreateAndRetrieveLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel",
		Description: "the description"}

	err := s.Annotations.CreateLabel(ctx, label)
	AssertNoError(t, err)

	retrievedLabel, err := s.Annotations.GetLabelById(ctx, label.Id.String())

	diff := deep.Equal(label, retrievedLabel)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical image structs, but got different fields: %v", diff))
	}

}

func TestDeleteLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}

	err := s.Annotations.CreateLabel(ctx, label)
	err = s.Annotations.DeleteLabel(ctx, label)

	AssertNoError(t, err)

	label, err = s.Annotations.GetLabelById(ctx, label.Id.String())
	AssertError(t, err)

}

func TestDeletingUsedLabelShouldFail(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}

	err := s.Annotations.CreateLabel(ctx, label)
	AssertNoError(t, err)
	image := &m.Image{Data: testImage}
	collection := &m.Collection{Name: "mycollection"}
	err = s.Collections.Create(ctx, collection)
	err = s.Images.Save(ctx, image, collection)
	err = s.Annotations.ApplyLabelToImage(ctx, label, image, collection)
	AssertNoError(t, err)

	err = s.Annotations.DeleteLabel(ctx, label)
	AssertError(t, err)

}

func TestDeleteLabeledImageAndItsAssociatedLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}
	image := &m.Image{Data: testImage}
	collection := &m.Collection{Name: "mycollection"}

	s.Annotations.CreateLabel(ctx, label)
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	s.Annotations.ApplyLabelToImage(ctx, label, image, collection)

	err := s.Images.Delete(ctx, image)
	err = s.Annotations.DeleteLabel(ctx, label)
	_, err = s.Annotations.GetLabelById(ctx, label.Id.String())
	AssertError(t, err)

}

func TestRemovingLabelFromImage(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: "mylabel"}
	collection := &m.Collection{Name: "myimageset"}
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	err := s.Annotations.CreateLabel(ctx, label)
	err = s.Annotations.ApplyLabelToImage(ctx, label, image, collection)
	err = s.Annotations.RemoveAnnotationFromImage(ctx, image.Annotations[0], image, collection)
	AssertNoError(t, err)

	if len(image.Annotations) != 0 {
		t.Fatalf("expected to retrieve image with 0 label, but got %v", len(image.Annotations))
	}

}
