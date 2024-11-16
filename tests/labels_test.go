package tests

import (
	m "go-image-annotator/models"
	"testing"
)

func TestInvalidLabelShouldFail(t *testing.T) {
	s, ctx := NewTestComponents(t)
	label := &m.Label{Name: "the name with spaces"}

	label, err := s.Labels.Create(ctx, label)
	AssertError(t, err)
}

func TestCreateLabel(t *testing.T) {
	s, ctx := NewTestComponents(t)
	label := &m.Label{Name: "thelabel",
		Description: "the description"}

	label, err := s.Labels.Create(ctx, label)
	AssertNoError(t, err)

	retrievedLabel, err := s.Labels.GetOne(ctx, label.Id.String())

	if label.Name != retrievedLabel.Name {
		t.Fatalf("expected to retrieve identical label names. Wanted %v, got %v", label.Name, retrievedLabel.Name)
	}

	if label.Description != retrievedLabel.Description {
		t.Fatalf("expected to retrieve identical label descriptions. Wanted %v, got %v", label.Description, retrievedLabel.Description)
	}

}
