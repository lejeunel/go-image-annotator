package tests

import (
	"context"

	a "datahub/app"
	an "datahub/app/annotator"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	e "datahub/errors"
	"testing"
)

func InitializeAuthTestsWithAdminEntitlement(t *testing.T) (*a.App, *loc.Site, *clc.Collection, *im.Image, *lbl.Label, context.Context) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "admin")
	ctx = context.WithValue(ctx, "groups", "mygroup")
	image, _ := im.New(testPNGImage)
	collection, _ := clc.New("mycollection", "", "mygroup")
	site, _ := loc.NewSite("thesite", "mygroup")
	label, _ := lbl.New("mylabel", "")
	s.Locations.SaveSite(ctx, site)
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	s.Labels.Create(ctx, label)

	return &s, site, collection, image, label, ctx

}

func TestCreatingSiteRequiresGroupMembership(t *testing.T) {
	s, _, _, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "im-contrib")

	site, _ := loc.NewSite("another-site", "another-group")
	err := s.Locations.SaveSite(ctx, site)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

}

func TestCreatingCollectionRequiresPermission(t *testing.T) {

	tests := map[string]struct {
		roles                   string
		wantEntitlementError    bool
		wantGroupOwnershipError bool
	}{
		"admin role should succeed":                           {roles: "admin", wantEntitlementError: false, wantGroupOwnershipError: false},
		"im-contrib role without group ownership should fail": {roles: "im-contrib", wantEntitlementError: false, wantGroupOwnershipError: true},
		"annotation-contrib role should fail":                 {roles: "annotation-contrib", wantEntitlementError: true, wantGroupOwnershipError: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _, _, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)
			ctx = context.WithValue(ctx, "entitlements", tc.roles)

			newCollection, _ := clc.New("my-new-collection", "", "")
			err := s.Collections.Create(ctx, newCollection)

			switch {
			case tc.wantEntitlementError:
				AssertErrorIs(t, err, e.ErrEntitlement)
			case tc.wantGroupOwnershipError:
				AssertErrorIs(t, err, e.ErrGroupOwnership)
			default:
				AssertNoError(t, err)
			}

		})
	}
}

func TestDeletingACollectionWithGroupOwnership(t *testing.T) {

	s, _, collection, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	collection, err := s.Collections.Find(ctx, collection.Id)
	AssertNoError(t, err)

	ctx = context.WithValue(ctx, "groups", "another-group")
	err = s.Images.DeleteCollection(ctx, collection)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

	ctx = context.WithValue(ctx, "groups", "mygroup")
	err = s.Images.DeleteCollection(ctx, collection)
	AssertNoError(t, err)
}

func TestUpdatingGroupOfCollectionRequiresToBeMemberOfThatGroup(t *testing.T) {

	s, _, collection, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "another-group")
	_, err := s.Collections.Update(ctx, collection.Name,
		clc.CollectionUpdatables{Name: "newname", Group: "another-group"})
	AssertErrorIs(t, err, e.ErrGroupOwnership)

}

func TestDeletingACollectionWithAdminEntitlementSucceeds(t *testing.T) {

	s, _, collection, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "another-group")
	err := s.Images.DeleteCollection(ctx, collection)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

	ctx = context.WithValue(ctx, "entitlements", "admin")
	err = s.Images.DeleteCollection(ctx, collection)
	AssertNoError(t, err)
}

func TestDeletingAnImageWithAdminEntitlementSucceeds(t *testing.T) {

	s, _, _, image, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "another-group")

	err := s.Images.RemoveFromCollection(ctx, image)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

	ctx = context.WithValue(ctx, "entitlements", "admin")
	err = s.Images.RemoveFromCollection(ctx, image)
	AssertNoError(t, err)
}

func TestDeletingAnImageWithGroupOwnershipSucceeds(t *testing.T) {

	s, _, _, image, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "another-group")

	err := s.Images.RemoveFromCollection(ctx, image)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

	ctx = context.WithValue(ctx, "groups", "mygroup")
	err = s.Images.RemoveFromCollection(ctx, image)
	AssertNoError(t, err)
}

func TestUpdatingImagesRequiresPermission(t *testing.T) {

	s, _, _, image, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "annotation-contrib")
	_, err := s.Images.Update(ctx, image.Id, im.ImageUpdatables{Site: "thesite",
		Camera: "", CapturedAt: "2006-01-02T15:04:05.000Z",
		Type_: "thermal"})
	AssertErrorIs(t, err, e.ErrEntitlement)
}

func TestCreatingLabelRequiresPermission(t *testing.T) {

	s, _, _, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	label, _ := lbl.New("newlabel", "")
	ctx = context.WithValue(ctx, "entitlements", "annotation-contrib")
	err := s.Labels.Create(ctx, label)
	AssertNoError(t, err)

	ctx = context.WithValue(ctx, "entitlements", "viewer")
	secondLabel, _ := lbl.New("anotherlabel", "")
	err = s.Labels.Create(ctx, secondLabel)
	AssertErrorIs(t, err, e.ErrEntitlement)

}

func TestDeletingLabelRequiresPermission(t *testing.T) {

	s, _, _, _, label, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "viewer")
	err := s.Labels.Delete(ctx, label)
	AssertErrorIs(t, err, e.ErrEntitlement)

	ctx = context.WithValue(ctx, "entitlements", "annotation-contrib")
	err = s.Labels.Delete(ctx, label)
	AssertNoError(t, err)

}

func TestApplyingBBoxRequiresEntitlement(t *testing.T) {

	s, _, _, image, label, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "viewer")
	bbox, _ := im.NewBoundingBox(10, 10, 11, 15)
	bbox.Annotate(label)
	err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertErrorIs(t, err, e.ErrEntitlement)

	ctx = context.WithValue(ctx, "entitlements", "annotation-contrib")
	err = s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertNoError(t, err)

}

func TestApplyingBBoxRequiresGroupOwnership(t *testing.T) {

	s, _, _, image, label, ctx := InitializeAuthTestsWithAdminEntitlement(t)

	ctx = context.WithValue(ctx, "entitlements", "annotation-contrib")
	ctx = context.WithValue(ctx, "groups", "another-group")
	bbox, _ := im.NewBoundingBox(10, 10, 11, 15)
	bbox.Annotate(label)
	err := s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

	ctx = context.WithValue(ctx, "groups", "mygroup")
	err = s.Images.Annotations.UpsertBoundingBox(ctx, bbox, image)
	AssertNoError(t, err)

}

func TestGeneratingAnnotatorStateDoesNotRequirePermission(t *testing.T) {

	s, _, _, image, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)
	annotator := an.NewAnnotator(s.Labels, s.Images,
		s.Collections, s.Locations, s.Authorizer, s.Logger, 640)

	ctx = context.WithValue(ctx, "entitlements", "")
	_, err := annotator.MakeState(ctx,
		an.AnnotatorRequest{ImageId: image.Id, CollectionId: image.Collection.Id})
	AssertNoError(t, err)

}

func TestAddingImageRequiresGroupOwnership(t *testing.T) {

	s, _, collection, _, _, ctx := InitializeAuthTestsWithAdminEntitlement(t)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "another-group")
	newImage, _ := im.New(testPNGImage)
	err := s.Images.Save(ctx, newImage, collection)
	AssertErrorIs(t, err, e.ErrGroupOwnership)

}
