package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"
	"testing"
)

func InitializeLocationsTests(t *testing.T) (*a.App, context.Context, *clc.Collection, *im.BaseImage) {
	s, _, ctx := a.NewTestApp(t, true)

	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("thename")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)
	base, _ := s.Images.GetBase(ctx, image.Id, im.FetchMetaOnly)

	return &s, ctx, collection, base

}

func TestCreateAndRetrieveSite(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	site, _ := loc.NewSite("thelocation")
	err := s.Locations.SaveSite(ctx, site)
	AssertNoError(t, err)

	retrievedSite, err := s.Locations.FindSite(ctx, site.Id)
	AssertNoError(t, err)
	if !retrievedSite.Id.Equal(&site.Id) {
		t.Fatalf("expected to retrieve site with id %v, but got %v", site.Id, retrievedSite.Id)
	}

}

func TestUpdateSiteName(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	site, _ := loc.NewSite("mysite")
	s.Locations.SaveSite(ctx, site)

	site.Name = "newname"
	_, err := s.Locations.UpdateSite(ctx, site)
	AssertNoError(t, err)

	retrieved, _ := s.Locations.FindSite(ctx, site.Id)

	AssertDeepEqual(t, retrieved, site, "site")
}

func TestUpdateSiteGroup(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	site, _ := loc.NewSite("mysite")
	s.Locations.SaveSite(ctx, site)

	site.Group = "anothergroup"
	_, err := s.Locations.UpdateSite(ctx, site)
	AssertNoError(t, err)

	retrieved, _ := s.Locations.FindSite(ctx, site.Id)

	AssertDeepEqual(t, retrieved, site, "site")
}

func TestCreateAndRetrieveSiteByName(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	siteName := "thesite"
	otherSiteName := "theothersite"
	site, _ := loc.NewSite(siteName)
	otherSite, _ := loc.NewSite(otherSiteName)
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveSite(ctx, otherSite)

	retrievedSite, err := s.Locations.FindSiteByName(ctx, otherSiteName)
	AssertNoError(t, err)
	AssertDeepEqual(t, *otherSite, *retrievedSite, "site")

}
func TestAssignCameraToSite(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	cameraName := "thecamera"
	site, _ := loc.NewSite("thelocation")
	camera, _ := loc.NewCamera(cameraName, site)
	s.Locations.SaveSite(ctx, site)
	err := s.Locations.SaveCamera(ctx, camera)
	AssertNoError(t, err)
	retrievedCamera, err := s.Locations.FindCamera(ctx, camera.Id)
	AssertDeepEqual(t, retrievedCamera.Site, site, "site")

}

func TestRetrieveCameraById(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	cameraName := "thecamera"
	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	camera, _ := loc.NewCamera(cameraName, site)
	err := s.Locations.SaveCamera(ctx, camera)

	retrievedCamera, err := s.Locations.FindCamera(ctx, camera.Id)
	AssertNoError(t, err)
	AssertDeepEqual(t, retrievedCamera, camera, "camera")
}

func TestRetrieveCameraByName(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	firstCamera, _ := loc.NewCamera("first-camera", site)
	s.Locations.SaveCamera(ctx, firstCamera)
	secondCamera, _ := loc.NewCamera("second-camera", site)
	s.Locations.SaveCamera(ctx, secondCamera)

	retrievedCamera, err := s.Locations.FindCameraByName(ctx,
		site, "first-camera")
	AssertNoError(t, err)
	if retrievedCamera.Name != "first-camera" {
		t.Fatalf("expected to retrieve camera with name %v, but got %v",
			"first-camera", retrievedCamera.Name)
	}
}

func TestCreateCamera(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	cameraName := "thecamera"
	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	camera, _ := loc.NewCamera(cameraName, site, loc.WithTransmitter("thetransmitter"))
	err := s.Locations.SaveCamera(ctx, camera)

	retrievedCamera, err := s.Locations.FindCamera(ctx,
		camera.Id)
	AssertNoError(t, err)
	if retrievedCamera.Transmitter != camera.Transmitter {
		t.Fatalf("expected to retrieve camera transmitter %v, but got %v",
			camera.Transmitter, retrievedCamera.Transmitter)
	}

	AssertDeepEqual(t, retrievedCamera, camera, "camera")
}

func TestCameraInheritsGroupFieldFromSite(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	cameraName := "thecamera"
	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	camera, _ := loc.NewCamera(cameraName, site)
	s.Locations.SaveCamera(ctx, camera)
	retrievedCamera, _ := s.Locations.FindCamera(ctx,
		camera.Id)
	if retrievedCamera.Group != site.Group {
		t.Fatalf("expected to retrieve camera group %v, but got %v",
			camera.Group, retrievedCamera.Group)
	}
}

func TestRetrieveImageWithTransmitter(t *testing.T) {
	s, ctx, collection, image := InitializeLocationsTests(t)
	cameraName := "thecamera"
	transmitterName := "thetransmitter"
	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	camera, _ := loc.NewCamera(cameraName, site, loc.WithTransmitter(transmitterName))
	s.Locations.SaveCamera(ctx, camera)
	s.Images.AssignCamera(ctx, camera.Id, image.Id)

	retrievedImage, _ := s.Images.Find(ctx, image.Id, collection.Id,
		im.FetchMetaOnly)

	if retrievedImage.GetTransmitter() != transmitterName {
		t.Fatalf("expected to retrieve transmitter %v, but got %v",
			transmitterName, retrievedImage.GetTransmitter())
	}

}

func TestCreatingInvalidSiteNameShouldFail(t *testing.T) {
	tests := map[string]struct {
		name string
	}{
		"with spaces":         {name: "the name with spaces"},
		"with capitals":       {name: "SiTe NaMe"},
		"with specials chars": {name: "s1T3&n4m3"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := loc.NewSite(tc.name)
			AssertErrorIs(t, err, e.ErrResourceName)
		})
	}
}

func TestCreatingInvalidCameraNameShouldFail(t *testing.T) {
	tests := map[string]struct {
		name string
	}{
		"with spaces":         {name: "the name with spaces"},
		"with capitals":       {name: "SiTe NaMe"},
		"with specials chars": {name: "s1T3&n4m3"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, _, ctx := a.NewTestApp(t, true)
			site, _ := loc.NewSite("thesite")
			s.Locations.SaveSite(ctx, site)

			_, err := loc.NewCamera(tc.name, site)
			AssertErrorIs(t, err, e.ErrResourceName)
		})
	}
}

func TestSiteDuplicateNameShouldFail(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	site, _ := loc.NewSite("my-site")
	newSite, _ := loc.NewSite("my-site")
	s.Locations.SaveSite(ctx, site)
	err := s.Locations.SaveSite(ctx, newSite)

	AssertErrorIs(t, err, e.ErrDuplication)
}

func TestCameraDuplicateNameShouldFail(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, true)

	site, _ := loc.NewSite("my-site")
	camera, _ := loc.NewCamera("my-camera", site)
	s.Locations.SaveSite(ctx, site)
	err := s.Locations.SaveCamera(ctx, camera)
	AssertNoError(t, err)

	err = s.Locations.SaveCamera(ctx, camera)
	AssertErrorIs(t, err, e.ErrDuplication)

}

func TestAssignNonExistingCameraToImageShouldFail(t *testing.T) {
	s, ctx, _, image := InitializeLocationsTests(t)

	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)

	nonExistingCamera, _ := loc.NewCamera("thecamera", site)

	err := s.Images.AssignCamera(ctx, nonExistingCamera.Id, image.Id)
	AssertErrorIs(t, err, e.ErrNotFound)

}

func TestCreateCameraWithNonExistingSiteToImageShouldFail(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	nonExistingSite, _ := loc.NewSite("thelocation")
	camera, _ := loc.NewCamera("thecamera", nonExistingSite)
	err := s.Locations.SaveCamera(ctx, camera)
	AssertErrorIs(t, err, e.ErrDependency)
}

func TestAssignCameraToImageShouldAlsoAddSite(t *testing.T) {
	s, ctx, collection, image := InitializeLocationsTests(t)

	site, _ := loc.NewSite("thelocation")
	camera, _ := loc.NewCamera("thecamera", site)
	s.Locations.SaveSite(ctx, site)

	s.Locations.SaveCamera(ctx, camera)

	err := s.Images.AssignCamera(ctx, camera.Id, image.Id)
	AssertNoError(t, err)

	retr, _ := s.Images.Find(ctx, image.Id, collection.Id, im.FetchMetaOnly)

	if retr.GetSiteName() == "" {
		t.Fatal("expected that image gets assigned a site, but it is empty")
	}
	if retr.GetCameraName() == "" {
		t.Fatal("expected that image gets assigned a camera, but it is empty")
	}

}

func TestDeleteSite(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	err := s.Locations.DeleteSite(ctx, site.Id)
	AssertNoError(t, err)
	_, err = s.Locations.FindSite(ctx, site.Id)
	AssertErrorIs(t, err, e.ErrNotFound)
}

func TestDeleteCamera(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	site, _ := loc.NewSite("thelocation")
	camera, _ := loc.NewCamera("thecamera", site)
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveCamera(ctx, camera)
	s.Locations.DeleteCamera(ctx, camera.Id)

	_, err := s.Locations.FindCamera(ctx, camera.Id)
	AssertErrorIs(t, err, e.ErrNotFound)
}

func TestUpdateSiteOfCameraShouldReassign(t *testing.T) {
	s, ctx, _, image := InitializeLocationsTests(t)

	firstSite, _ := loc.NewSite("thelocation")
	secondSite, _ := loc.NewSite("thesecondlocation")
	camera, _ := loc.NewCamera("thecamera", firstSite)
	s.Locations.SaveSite(ctx, firstSite)
	s.Locations.SaveSite(ctx, secondSite)
	s.Locations.SaveCamera(ctx, camera)
	updatedCamera, err := s.Locations.UpdateCamera(ctx, camera.Id,
		loc.CameraUpdatables{Name: "newname", SiteName: "thesecondlocation"})
	AssertNoError(t, err)

	s.Images.AssignCamera(ctx, camera.Id, image.Id)

	if updatedCamera.Site.Name != "thesecondlocation" {
		t.Fatalf("expected to update site of camera to %v, but got %v",
			"thesecondlocation", updatedCamera.Name)
	}

	retrievedImages, _, _ := s.Images.List(ctx, *im.NewImageFilter(im.WithCameraId(camera.Id)),
		im.OrderingArgs{}, g.PaginationParams{Page: 1, PageSize: 2}, im.FetchMetaOnly)
	if retrievedImages[0].GetSiteName() != secondSite.Name {
		t.Fatal("expected to update site of image but did not")
	}
}

func TestUpdateCameraName(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	site, _ := loc.NewSite("thelocation")
	camera, _ := loc.NewCamera("thecamera", site)
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveCamera(ctx, camera)
	updatedCamera, err := s.Locations.UpdateCamera(ctx, camera.Id,
		loc.CameraUpdatables{Name: "newname", SiteName: "thelocation"})
	AssertNoError(t, err)

	if updatedCamera.Name != "newname" {
		t.Fatalf("expected to update name of camera to %v, but got %v",
			"newname", updatedCamera.Name)
	}
}

func TestUpdateCameraTransmitter(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	site, _ := loc.NewSite("thelocation")
	camera, _ := loc.NewCamera("thecamera", site, loc.WithTransmitter("thetransmitter"))
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveCamera(ctx, camera)
	newTransmitterName := "thenewtransmitter"
	updatedCamera, err := s.Locations.UpdateCamera(ctx, camera.Id,
		loc.CameraUpdatables{Name: "newname", SiteName: "thelocation",
			Transmitter: newTransmitterName})
	AssertNoError(t, err)

	if updatedCamera.Transmitter != newTransmitterName {
		t.Fatalf("expected to update name of camera to %v, but got %v",
			newTransmitterName, updatedCamera.Name)
	}
}

func TestReassignCameraToAnotherSite(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	firstSite, _ := loc.NewSite("firstsite")
	secondSite, _ := loc.NewSite("secondsite")
	s.Locations.SaveSite(ctx, firstSite)
	camera, _ := loc.NewCamera("thecamera", firstSite)
	s.Locations.SaveSite(ctx, secondSite)
	s.Locations.SaveCamera(ctx, camera)
	updatedCamera, err := s.Locations.UpdateCamera(ctx, camera.Id,
		loc.CameraUpdatables{Name: "thecamera", SiteName: "secondsite"})
	AssertNoError(t, err)

	if updatedCamera.Site.Name != "secondsite" {
		t.Fatalf("expected to update site of camera to %v, but got %v",
			"secondsite", updatedCamera.Site.Name)
	}
}

func TestReassigningCameraToSiteWithAlreadyExistingNameShouldFail(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	firstSite, _ := loc.NewSite("firstsite")
	secondSite, _ := loc.NewSite("secondsite")
	s.Locations.SaveSite(ctx, firstSite)
	s.Locations.SaveSite(ctx, secondSite)
	firstCamera, _ := loc.NewCamera("thecamera", firstSite)
	secondCamera, _ := loc.NewCamera("thecamera", secondSite)
	s.Locations.SaveCamera(ctx, firstCamera)
	s.Locations.SaveCamera(ctx, secondCamera)
	_, err := s.Locations.UpdateCamera(ctx, secondCamera.Id,
		loc.CameraUpdatables{Name: "thecamera", SiteName: "firstsite"})
	AssertErrorIs(t, err, e.ErrDuplication)

}

func TestListImageOfCamera(t *testing.T) {
	s, ctx, _, _ := InitializeLocationsTests(t)

	firstSite, _ := loc.NewSite("firstsite")
	firstCamera, _ := loc.NewCamera("firstcamera", firstSite)
	secondSite, _ := loc.NewSite("secondsite")
	secondCamera, _ := loc.NewCamera("secondcamera", secondSite)
	s.Locations.SaveSite(ctx, firstSite)
	s.Locations.SaveSite(ctx, secondSite)
	s.Locations.SaveCamera(ctx, firstCamera)
	s.Locations.SaveCamera(ctx, secondCamera)

	firstImage, _ := im.New(testJPGImage)
	secondImage, _ := im.New(testJPGImage)
	collection, _ := clc.New("thecollection")
	s.Collections.Create(ctx, collection)

	s.Images.Save(ctx, firstImage, collection)
	s.Images.AssignCamera(ctx, firstCamera.Id, firstImage.Id)

	s.Images.Save(ctx, secondImage, collection)
	s.Images.AssignCamera(ctx, secondCamera.Id, secondImage.Id)

	retrieved, _, _ := s.Images.List(ctx, im.FilterArgs{CameraId: &secondCamera.Id},
		*im.NewImageDefaultOrderingArgs(), g.PaginationParams{Page: 1, PageSize: 2},
		im.FetchMetaOnly)

	if len(retrieved) != 1 {
		t.Fatalf("expected to retrieve 1 image, but got %v", len(retrieved))
	}
	if retrieved[0].GetCameraName() != "secondcamera" {
		t.Fatalf("expected to retrieve image with camera %v , got %v", "secondcamera", retrieved[0].GetCameraName())
	}

}

func InitCollection(t *testing.T, a *a.App, ctx context.Context,
	image *im.Image, collectionName, siteName, cameraName string) {

	collection, _ := clc.New(collectionName, clc.WithGroup("a-group"))
	a.Collections.Create(ctx, collection)
	a.Images.Save(ctx, image, collection)
	site, _ := loc.NewSite(siteName, loc.WithGroupOption("a-group"))
	camera, _ := loc.NewCamera(cameraName, site)
	a.Locations.SaveSite(ctx, site)
	a.Locations.SaveCamera(ctx, camera)
	a.Images.AssignCamera(ctx, camera.Id, image.Id)

}

func TestFilterSitesByCollection(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "a-group")

	image, _ := im.New(testJPGImage)
	InitCollection(t, &s, ctx, image, "my-collection", "my-site", "my-camera")
	InitCollection(t, &s, ctx, image, "other-collection", "other-site", "other-camera")

	retrieved, pagination, _ := s.Locations.List(ctx, *loc.NewSiteFilter(loc.WithCollection("my-collection")),
		loc.SiteAlphabeticalOrdering,
		g.PaginationParams{Page: 1, PageSize: 2})

	if len(retrieved) != 1 {
		t.Fatalf("expected to retrieve one site but got %v", len(retrieved))
	}

	if pagination.TotalRecords != 1 {
		t.Fatalf("expected to retrieve one site but got %v", pagination.TotalRecords)
	}

}

func TestFilterSitesByGroup(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "a-group")

	site, _ := loc.NewSite("my-site", loc.WithGroupOption("a-group"))
	s.Locations.SaveSite(ctx, site)
	otherSite, _ := loc.NewSite("my-other-site")
	s.Locations.SaveSite(ctx, otherSite)

	retrieved, _, err := s.Locations.List(ctx, *loc.NewSiteFilter(loc.WithGroup("a-group")), loc.SiteAlphabeticalOrdering, g.OneItemPaginationParams)
	AssertNoError(t, err)

	if len(retrieved) == 0 {
		t.Fatal("expected to retrieve one site but got none")
	}

	if retrieved[0].Name != "my-site" {
		t.Fatalf("expected to retrieve site with name %v but got %v", "first-site", retrieved[0].Name)
	}
}

func TestSitesAreReturnedInAlphabeticalOrder(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "thegroup")
	nonOrderedNames := []string{"c", "b", "a"}
	for _, name := range nonOrderedNames {
		site, _ := loc.NewSite(name, loc.WithGroupOption("thegroup"))
		err := s.Locations.SaveSite(ctx, site)
		AssertNoError(t, err)
	}

	orderedNames := []string{"a", "b", "c"}
	retrievedSites, _, err := s.Locations.List(ctx, loc.FilterArgs{},
		loc.SiteAlphabeticalOrdering,
		g.PaginationParams{Page: 1, PageSize: 3})
	AssertNoError(t, err)
	for i, site := range retrievedSites {
		if site.Name != orderedNames[i] {
			t.Fatalf("expected to retrieve site name %v, but got %v", orderedNames[i], site.Name)
		}
	}

}

func TestPatchCamera(t *testing.T) {

	tests := []struct {
		testName   string
		updateName bool
		updateSite bool
	}{
		{testName: "update name", updateName: true, updateSite: false},
		{testName: "update site", updateName: false, updateSite: true},
		{testName: "update both", updateName: true, updateSite: true},
	}

	for _, tc := range tests {
		s, _, ctx := a.NewTestApp(t, true)

		image, _ := im.New(testJPGImage)
		collection, _ := clc.New("my-collection")
		s.Collections.Create(ctx, collection)
		s.Images.Save(ctx, image, collection)

		site, _ := loc.NewSite("my-site")
		newSite, _ := loc.NewSite("my-new-site")
		camera, _ := loc.NewCamera("my-camera", site)
		s.Locations.SaveSite(ctx, site)
		s.Locations.SaveSite(ctx, newSite)
		s.Locations.SaveCamera(ctx, camera)
		s.Images.AssignCamera(ctx, camera.Id, image.Id)

		patches := []g.JSONPatch{}
		if tc.updateName {
			patches = append(patches,
				g.JSONPatch{Operation: "replace", Path: "/name", Value: "new-camera-name"})
			_, err := s.Locations.PatchCamera(ctx, camera.Id, patches)
			AssertNoError(t, err)
			retrievedCamera, _ := s.Locations.FindCamera(ctx, camera.Id)
			if retrievedCamera.Name != "new-camera-name" {
				t.Fatalf("expected to rename camera to %v but did not. Got %v instead",
					"new-camera-name", retrievedCamera.Name)
			}
		}
		if tc.updateSite {
			patches = append(patches,
				g.JSONPatch{Operation: "replace", Path: "/site", Value: "my-new-site"})
			_, err := s.Locations.PatchCamera(ctx, camera.Id, patches)
			AssertNoError(t, err)
			retrievedCamera, _ := s.Locations.FindCamera(ctx, camera.Id)
			if retrievedCamera.Site.Name != "my-new-site" {
				t.Fatalf("expected to change site to %v, but did not. Got %v instead",
					"my-new-site", retrievedCamera.Site.Name)
			}
		}

	}

}

func TestListSitesByGroup(t *testing.T) {
	s, _, ctx := a.NewTestApp(t, false)
	ctx = context.WithValue(ctx, "entitlements", "im-contrib")
	ctx = context.WithValue(ctx, "groups", "first-group|second-group")
	firstSite, _ := loc.NewSite("first-site", loc.WithGroupOption("first-group"))
	secondSite, _ := loc.NewSite("second-site", loc.WithGroupOption("second-group"))
	s.Locations.SaveSite(ctx, firstSite)
	s.Locations.SaveSite(ctx, secondSite)

	retrievedSites, paginationMeta, err := s.Locations.List(ctx, loc.NewSiteByGroupFilter("first-group"),
		loc.SiteAlphabeticalOrdering,
		g.PaginationParams{Page: 1, PageSize: 3})
	AssertNoError(t, err)
	if len(retrievedSites) != 1 {
		t.Fatalf("expected to retrieve one site, but got %v",
			len(retrievedSites))
	}
	if paginationMeta.TotalRecords != 1 {
		t.Fatalf("expected to retrieve pagination with one record, but got %v",
			paginationMeta.TotalRecords)
	}

}

func TestClearCameraFromImage(t *testing.T) {
	s, ctx, collection, image := InitializeLocationsTests(t)
	cameraName := "thecamera"
	site, _ := loc.NewSite("thelocation")
	s.Locations.SaveSite(ctx, site)
	camera, _ := loc.NewCamera(cameraName, site)
	s.Locations.SaveCamera(ctx, camera)
	s.Images.AssignCamera(ctx, camera.Id, image.Id)

	err := s.Images.UnassignCamera(ctx, image.Id)
	AssertNoError(t, err)

	retrievedImage, _ := s.Images.Find(ctx, image.Id, collection.Id,
		im.FetchMetaOnly)

	if retrievedImage.Camera != nil {
		t.Fatal("expected to retrieve image with nil camera")
	}

}
