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

func TestLabelDuplicateNameShouldFail(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	label := &m.Label{Name: "mylabel"}
	newLabel := &m.Label{Name: "mylabel"}
	s.Annotations.CreateLabel(ctx, label)
	err := s.Annotations.CreateLabel(ctx, newLabel)

	AssertError(t, err)
}

func TestCreateAndRetrieveLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel",
		Description: "the description"}

	err := s.Annotations.CreateLabel(ctx, label)
	AssertNoError(t, err)

	retrievedLabel, err := s.Annotations.GetOrdinal(ctx, 0)
	AssertNoError(t, err)
	if retrievedLabel == nil {
		t.Fatal("expected to retrieve 1 label, but got none")
	}

	diff := deep.Equal(*label, *retrievedLabel)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical label structs, but got different fields: %v", diff))
	}

}

func TestDeleteLabel(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}

	err := s.Annotations.CreateLabel(ctx, label)
	err = s.Annotations.DeleteLabel(ctx, label)

	AssertNoError(t, err)

	retrievedLabel, _ := s.Annotations.GetOrdinal(ctx, 0)
	if retrievedLabel != nil {
		t.Fatal("expected to retrieve 0 labels, but got one")
	}

}

func TestDeletingUsedLabelShouldFail(t *testing.T) {
	s, ctx := NewTestApp(t, 2)
	label := &m.Label{Name: "thelabel"}

	s.Annotations.CreateLabel(ctx, label)
	image := &m.Image{Data: testImage}
	collection := &m.Collection{Name: "mycollection"}
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	err := s.Annotations.ApplyLabelToImage(ctx, label, image, collection)
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

	s.Images.Delete(ctx, image, collection)
	s.Annotations.DeleteLabel(ctx, label)
	retrievedLabel, _ := s.Annotations.GetOrdinal(ctx, 0)
	if retrievedLabel != nil {
		t.Fatal("expected to retrieve 0 labels, but got one")
	}

}

func TestRemovingLabelFromImage(t *testing.T) {
	s, ctx := NewTestApp(t, 2)

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: "mylabel"}
	collection := &m.Collection{Name: "myimageset"}
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	s.Annotations.CreateLabel(ctx, label)
	s.Annotations.ApplyLabelToImage(ctx, label, image, collection)
	s.Annotations.RemoveAnnotationFromImage(ctx, image.Annotations[0], image, collection)

	if len(image.Annotations) != 0 {
		t.Fatalf("expected to retrieve image with 0 label, but got %v", len(image.Annotations))
	}

}
