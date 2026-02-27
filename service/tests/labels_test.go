package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	"testing"
)

func InitLabelTest(t *testing.T) (*a.App, *clc.Collection, *im.Image, *lbl.Label, context.Context) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "admin")
	ctx = context.WithValue(ctx, "groups", "mygroup")
	image, _ := im.New(testPNGImage)
	collection, _ := clc.New("myimageset", clc.WithGroup("mygroup"))
	label, _ := lbl.New("mylabel", "mydescription")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	err := s.Labels.Create(ctx, label)
	AssertNoError(t, err)

	return &s, collection, image, label, ctx

}

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
			_, err := lbl.New(tc.name, "")
			AssertErrorIs(t, err, e.ErrResourceName)
		})
	}
}
func TestUpdateLabelDescription(t *testing.T) {
	s, _, _, _, ctx := InitLabelTest(t)

	label, _ := lbl.New("mynewlabel", "")
	s.Labels.Create(ctx, label)
	newDescription := "new description"
	err := s.Labels.Update(ctx, label, lbl.Updatables{Description: newDescription})
	AssertNoError(t, err)

	if label.Description != newDescription {
		t.Fatalf("expected to store updated description in label with value %v but got %v",
			newDescription, label.Description)
	}

	retrieved, _ := s.Labels.Find(ctx, label.Id)
	if retrieved.Description != newDescription {
		t.Fatalf("expected to retrieve label with updated description %v but got %v",
			newDescription, retrieved.Description)
	}

}

func TestLabelDuplicateNameShouldFail(t *testing.T) {
	s, _, _, _, ctx := InitLabelTest(t)

	newLabel, _ := lbl.New("mynewlabel", "")
	s.Labels.Create(ctx, newLabel)
	labelWithSameName, _ := lbl.New("mynewlabel", "")
	err := s.Labels.Create(ctx, labelWithSameName)

	AssertErrorIs(t, err, e.ErrDuplication)
}

func TestCreateAndRetrieveLabel(t *testing.T) {
	s, _, _, label, ctx := InitLabelTest(t)
	retrievedLabel, err := s.Labels.Find(ctx, label.Id)
	AssertNoError(t, err)
	if retrievedLabel == nil {
		t.Fatal("expected to retrieve 1 label, but got none")
	}

	AssertDeepEqual(t, *label, *retrievedLabel, "label")

}

func TestDeleteLabel(t *testing.T) {
	s, _, _, label, ctx := InitLabelTest(t)
	err := s.Labels.Delete(ctx, label)
	AssertNoError(t, err)

	res, _ := s.Labels.Find(ctx, label.Id)
	if res != nil {
		t.Fatal("expected to retrieve 0 labels, but got one")
	}

}

func TestDeletingUsedLabelShouldFail(t *testing.T) {
	s, _, image, label, ctx := InitLabelTest(t)
	err := s.Images.Annotations.ApplyLabel(ctx, label, image)
	AssertNoError(t, err)

	err = s.Labels.Delete(ctx, label)
	AssertErrorIs(t, err, e.ErrDependency)
}

func TestCreateAndRetrieveLabelHierarchy(t *testing.T) {
	s, _, _, _, ctx := InitLabelTest(t)
	parent, _ := lbl.New("theparent", "")
	child, _ := lbl.New("thechild", "")
	grandChild, _ := lbl.New("thegrandchlid", "")

	s.Labels.Create(ctx, parent)
	s.Labels.Create(ctx, child)
	s.Labels.Create(ctx, grandChild)
	err := s.Labels.SetParenting(ctx, child, parent)
	AssertNoError(t, err)
	err = s.Labels.SetParenting(ctx, grandChild, child)
	AssertNoError(t, err)

	retrievedChild, err := s.Labels.Find(ctx, child.Id)
	AssertNoError(t, err)
	retrievedGrandChild, err := s.Labels.Find(ctx, grandChild.Id)
	AssertNoError(t, err)

	if retrievedChild.Parent == nil {
		t.Fatalf("expected to retrieve label with a parent, but got none")
	}

	AssertDeepEqual(t, *retrievedChild.Parent, *parent, "child label")
	AssertDeepEqual(t, *retrievedGrandChild.Parent, *child, "grand-child label")
	AssertDeepEqual(t, *retrievedGrandChild.Parent.Parent, *parent, "parent label")

}

func TestFetchLabelByName(t *testing.T) {
	s, _, _, label, ctx := InitLabelTest(t)

	res, err := s.Labels.FindByName(ctx, label.Name)
	AssertNoError(t, err)
	AssertDeepEqual(t, label, res, "label")

}

func TestLabelsAreReturnedInAlphabeticalOrder(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "annotation-contrib")
	nonOrderedNames := []string{"c", "b", "a"}
	for _, name := range nonOrderedNames {
		label, _ := lbl.New(name, "")
		err := s.Labels.Create(ctx, label)
		AssertNoError(t, err)
	}

	orderedNames := []string{"a", "b", "c"}
	retrievedLabels, err := s.Labels.GetAllLabels(ctx)
	AssertNoError(t, err)
	for i, retrievedLabel := range retrievedLabels {
		if retrievedLabel.Name != orderedNames[i] {
			t.Fatalf("expected to retrieve label name %v, but got %v", orderedNames[i], retrievedLabel.Name)
		}
	}

}

func TestSerializeLabelHierarchyToCsv(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)
	parentLabel, _ := lbl.New("theparent", "")
	childLabel, _ := lbl.New("thechild", "")
	grandChildLabel, _ := lbl.New("thegrandchild", "")

	err := s.Labels.Create(ctx, parentLabel)
	AssertNoError(t, err)
	s.Labels.Create(ctx, childLabel)
	err = s.Labels.Create(ctx, grandChildLabel)
	s.Labels.SetParenting(ctx, childLabel, parentLabel)
	s.Labels.SetParenting(ctx, grandChildLabel, childLabel)

	expected := "thegrandchild,thechild,theparent"
	got := grandChildLabel.String()
	if got != expected {
		t.Fatalf("expected to retrieve serialized label hierarchy %v but got %v",
			expected, got)
	}

}
