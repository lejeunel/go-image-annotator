package tests

import (
	"context"
	a "datahub/app"
	an "datahub/app/annotator"
	pro "datahub/domain/annotation_profiles"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	g "datahub/generic"
	"testing"
)

func InitializeAnnotationProfileTests(t *testing.T) (*a.App, *clc.Collection, *im.Image, *lbl.Label, context.Context) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "admin")
	ctx = context.WithValue(ctx, "groups", "mygroup")
	image, _ := im.New(testPNGImage)
	collection, _ := clc.New("myimageset", clc.WithGroup("mygroup"))
	label, _ := lbl.New("my-label", "mydescription")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	s.Labels.Create(ctx, label)

	return &s, collection, image, label, ctx
}

func TestListingAnnotationProfiles(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	firstProfile := pro.New("first-profile")
	s.Profiles.Save(ctx, firstProfile)
	secondProfile := pro.New("second-profile")
	s.Profiles.Save(ctx, secondProfile)

	profiles, _, _ := s.Profiles.List(ctx, g.PaginationParams{Page: 1, PageSize: 1})
	if len(profiles) != 1 {
		t.Fatalf("expected page with one item, but got %v", len(profiles))
		if profiles[0].Name != "first-profile" {
			t.Fatalf("expected to retrieve profile with name %v but got %v",
				"first-profile", profiles[0].Name)
		}
		if profiles[1].Name != "second-profile" {
			t.Fatalf("expected to retrieve profile with name %v but got %v",
				"second-profile", profiles[1].Name)
		}
	}

}

func TestCreateAnnotationProfile(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	err := s.Profiles.Save(ctx, profile)
	AssertNoError(t, err)

	retrievedProfile, err := s.Profiles.Find(ctx, profile.Id)
	AssertNoError(t, err)
	AssertDeepEqual(t, retrievedProfile, profile, "annotation profile")
}

func TestFetchAnnotationProfileByName(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	retrievedProfile, err := s.Profiles.FindByName(ctx, profile.Name)
	AssertNoError(t, err)
	AssertDeepEqual(t, retrievedProfile, profile, "annotation profile")
}

func TestFetchNonExistingAnnotationProfileByName(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	_, err := s.Profiles.FindByName(ctx, "non-existing-profile")
	AssertErrorIs(t, err, e.ErrNotFound)
}

func TestDeleteAnnotationProfile(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)
	err := s.Profiles.Delete(ctx, profile)
	AssertNoError(t, err)

	_, err = s.Profiles.Find(ctx, profile.Id)
	AssertErrorIs(t, err, e.ErrNotFound)
}

func TestAddingLabelToAnnotationProfile(t *testing.T) {
	s, _, _, label, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	err := s.Profiles.AddLabel(ctx, profile, label)
	AssertNoError(t, err)

	retrievedProfile, err := s.Profiles.Find(ctx, profile.Id)
	if len(retrievedProfile.Labels) != 1 {
		t.Fatalf("expected to retrieve profile with one label, but got %v",
			len(retrievedProfile.Labels))
	}
}

func TestUpdateAnnotationProfileName(t *testing.T) {
	s, _, _, label, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)
	s.Profiles.AddLabel(ctx, profile, label)

	newProfileName := "new-name"
	updatedProfile, err := s.Profiles.Update(ctx, profile.Id,
		pro.ProfileUpdatables{Name: newProfileName})
	AssertNoError(t, err)

	retrievedUpdatedProfile, _ := s.Profiles.Find(ctx, profile.Id)
	AssertDeepEqual(t, updatedProfile, retrievedUpdatedProfile, "profile")

	if retrievedUpdatedProfile.Name != newProfileName {
		t.Fatalf("expected to retrieve profile with name %v, but got %v",
			"new-name", retrievedUpdatedProfile.Name)
	}
}

func TestAddingLabelSetToAnnotationProfile(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	label0, _ := lbl.New("first-label", "")
	label1, _ := lbl.New("second-label", "")
	s.Labels.Create(ctx, label0)
	s.Labels.Create(ctx, label1)

	err := s.Profiles.AddLabelSet(ctx, profile, []string{label0.Name, label1.Name})
	AssertNoError(t, err)

	retrievedProfile, err := s.Profiles.Find(ctx, profile.Id)
	if len(retrievedProfile.Labels) != 2 {
		t.Fatalf("expected to retrieve profile with two labels, but got %v",
			len(retrievedProfile.Labels))
	}
}

func TestAnnotationProfileLabelsAreInAlphabeticalOrder(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	labelA, _ := lbl.New("label-a", "")
	labelB, _ := lbl.New("label-b", "")
	s.Labels.Create(ctx, labelA)
	s.Labels.Create(ctx, labelB)

	err := s.Profiles.AddLabelSet(ctx, profile, []string{labelB.Name, labelA.Name})
	AssertNoError(t, err)

	retrievedProfile, err := s.Profiles.Find(ctx, profile.Id)
	if retrievedProfile.Labels[0].Name != "label-a" {
		t.Fatalf("expected to retrieve first label with name %v but got %v",
			"label-a",
			retrievedProfile.Labels[0].Name)
	}
}

func TestAddingSameLabelTwiceToProfileShouldFail(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	label0, _ := lbl.New("first-label", "")
	s.Labels.Create(ctx, label0)

	err := s.Profiles.AddLabelSet(ctx, profile, []string{label0.Name, label0.Name})
	AssertError(t, err)

}

func TestRemovingAllLabelsFromAnnotationProfile(t *testing.T) {
	s, _, _, label, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	err := s.Profiles.RemoveLabel(ctx, profile, label)
	AssertNoError(t, err)

	otherLabel, _ := lbl.New("other-label", "")
	s.Labels.Create(ctx, otherLabel)

	err = s.Profiles.ClearLabels(ctx, profile)
	if len(profile.Labels) != 0 {
		t.Fatalf("expected to get profile with no label, but got %v",
			len(profile.Labels))
	}
	AssertNoError(t, err)

	retrievedProfile, err := s.Profiles.Find(ctx, profile.Id)
	if len(retrievedProfile.Labels) != 0 {
		t.Fatalf("expected to retrieve profile with no label, but got %v",
			len(retrievedProfile.Labels))
	}
}

func TestRemovingLabelFromAnnotationProfile(t *testing.T) {
	s, _, _, label, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)

	err := s.Profiles.RemoveLabel(ctx, profile, label)
	AssertNoError(t, err)

	retrievedProfile, err := s.Profiles.Find(ctx, profile.Id)
	if len(retrievedProfile.Labels) != 0 {
		t.Fatalf("expected to retrieve profile with no label, but got %v",
			len(retrievedProfile.Labels))
	}
}

func TestAssignAnnotationProfile(t *testing.T) {
	s, collection, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)
	err := s.Collections.AssignProfile(ctx, profile, collection)
	AssertNoError(t, err)

	retrievedCollection, err := s.Collections.Find(ctx, collection.Id)
	AssertNoError(t, err)
	if retrievedCollection.Profile == nil {
		t.Fatal("expected to retrieve profile associated to collection when retriving by id, but got none")
	}
}

func TestAssignAnnotationProfileRetrieveByName(t *testing.T) {
	s, collection, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)
	err := s.Collections.AssignProfile(ctx, profile, collection)
	AssertNoError(t, err)

	retrievedCollection, err := s.Collections.FindByName(ctx, collection.Name)
	AssertNoError(t, err)
	if retrievedCollection.Profile == nil {
		t.Fatal("expected to retrieve profile associated to collection when retrieving by name, but got none")
	}
}

func TestUnAssignAnnotationProfileFromCollection(t *testing.T) {
	s, collection, _, _, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)
	s.Collections.AssignProfile(ctx, profile, collection)

	err := s.Collections.UnassignProfile(ctx, collection)
	AssertNoError(t, err)

	retrievedCollection, err := s.Collections.Find(ctx, collection.Id)
	if retrievedCollection.Profile != nil {
		t.Fatal("expected to retrieve empty profile, but got one")
	}

}

func TestFetchingLabelsInCollectionWithProfile(t *testing.T) {
	s, collection, _, _, ctx := InitializeAnnotationProfileTests(t)

	newLabel, _ := lbl.New("new-label", "")
	s.Labels.Create(ctx, newLabel)

	profile := pro.New("my-profile")
	s.Profiles.Save(ctx, profile)
	s.Profiles.AddLabel(ctx, profile, newLabel)
	s.Collections.AssignProfile(ctx, profile, collection)

	labels, err := s.Collections.GetAvailableLabels(ctx, collection)
	AssertNoError(t, err)
	if len(labels) != 1 {
		t.Fatalf("expected to retrieve set of available labels of size %v, but got %v",
			1, len(labels))
	}
}

func TestFetchingLabelsInCollectionWithoutProfileShouldReturnAllOfThem(t *testing.T) {
	s, collection, _, _, ctx := InitializeAnnotationProfileTests(t)

	labels, err := s.Collections.GetAvailableLabels(ctx, collection)
	AssertNoError(t, err)

	if len(labels) != 1 {
		t.Fatalf("expected to retrieve set of available labels of size %v, but got %v",
			1, len(labels))
	}
}

func TestAnnotatorWithAssignedProfile(t *testing.T) {
	s, collection, image, _, ctx := InitializeAnnotationProfileTests(t)
	annotator := an.NewAnnotator(s.Labels, s.Images,
		s.Collections, s.Locations, s.Authorizer, s.Logger, 640)

	profile := pro.New("my-profile")
	newLabelName := "new-label"
	s.Profiles.Save(ctx, profile)
	newLabel, _ := lbl.New(newLabelName, "")
	s.Labels.Create(ctx, newLabel)
	s.Profiles.AddLabel(ctx, profile, newLabel)
	s.Collections.AssignProfile(ctx, profile, collection)
	state, err := annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.Collection.Id})
	AssertNoError(t, err)

	availableLabels := state.AvailableLabels
	if len(availableLabels) != 1 {
		t.Fatalf("expected to retrieve one available label but got %v",
			len(availableLabels))
	}

	if availableLabels[0] != newLabelName {
		t.Fatalf("expected to retrieve available label named %v but got %v",
			newLabelName, len(availableLabels))
	}

}

func TestAnnotatingImageWithAssignedProfile(t *testing.T) {
	s, collection, image, forbiddenLabel, ctx := InitializeAnnotationProfileTests(t)

	profile := pro.New("my-profile")
	availableLabel := "new-label"
	s.Profiles.Save(ctx, profile)
	newLabel, _ := lbl.New(availableLabel, "")
	s.Labels.Create(ctx, newLabel)
	s.Profiles.AddLabel(ctx, profile, newLabel)
	s.Collections.AssignProfile(ctx, profile, collection)

	bbox, _ := im.NewBoundingBox(10, 10, 11, 15)
	bbox.Annotate(forbiddenLabel)
	err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertErrorIs(t, err, e.ErrForbiddenLabel)

	bbox.Annotate(newLabel)
	err = s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertNoError(t, err)

}

func TestListingAnnotationProfilesInAlphanumericalOrder(t *testing.T) {
	s, _, _, _, ctx := InitializeAnnotationProfileTests(t)

	aProfile := pro.New("a-profile")
	bProfile := pro.New("b-profile")
	s.Profiles.Save(ctx, bProfile)
	s.Profiles.Save(ctx, aProfile)

	profiles, _, _ := s.Profiles.List(ctx, g.PaginationParams{Page: 1, PageSize: 2})
	if profiles[0].Name != "a-profile" {
		t.Fatalf("expected to retrieve profile with name %v but got %v",
			"a-profile", profiles[0].Name)
	}

}
