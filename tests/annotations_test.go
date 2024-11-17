package tests

import (
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
			s, ctx := NewTestComponents(t)
			label := &m.Label{Name: tc.name}
			label, err := s.Annotations.Create(ctx, label)
			AssertError(t, err)
		})
	}
}

func TestCreateAndRetrieveLabel(t *testing.T) {
	s, ctx := NewTestComponents(t)
	label := &m.Label{Name: "thelabel",
		Description: "the description"}

	label, err := s.Annotations.Create(ctx, label)
	AssertNoError(t, err)

	retrievedLabel, err := s.Annotations.GetOne(ctx, label.Id.String())

	if label.Name != retrievedLabel.Name {
		t.Fatalf("expected to retrieve identical label names. Wanted %v, got %v", label.Name, retrievedLabel.Name)
	}

	if label.Description != retrievedLabel.Description {
		t.Fatalf("expected to retrieve identical label descriptions. Wanted %v, got %v", label.Description, retrievedLabel.Description)
	}

}

func TestDeleteLabel(t *testing.T) {
	s, ctx := NewTestComponents(t)
	label := &m.Label{Name: "thelabel"}

	label, err := s.Annotations.Create(ctx, label)
	err = s.Annotations.Delete(ctx, label)

	AssertNoError(t, err)

	label, err = s.Annotations.GetOne(ctx, label.Id.String())
	AssertError(t, err)

}

func TestDeletingUsedLabelShouldFail(t *testing.T) {
	s, ctx := NewTestComponents(t)
	label := &m.Label{Name: "thelabel"}

	label, _ = s.Annotations.Create(ctx, label)
	image := &m.Image{Data: testImage}
	image, _ = s.Annotations.ApplyLabelToImage(ctx, label, image)

	err := s.Annotations.Delete(ctx, label)

	AssertError(t, err)

}

func TestDeleteLabeledImageAndItsAssociatedLabel(t *testing.T) {
	s, ctx := NewTestComponents(t)
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
