package tests

import (
	"context"
	locpck "datahub/app/locationpicker"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"

	a "datahub/app"
	"testing"
)

func InitializeLocationPickerTests(t *testing.T) (context.Context, *a.App, *locpck.LocationPicker, *im.Image, *loc.Site, []*loc.Camera) {
	s, _, ctx := a.NewTestApp(t, true)
	label, _ := lbl.New("thelabel", "")
	s.Labels.Create(ctx, label)
	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("thecollection", "", "mygroup")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	ctx = context.WithValue(ctx, "groups", "mygroup")

	picker := locpck.NewLocationPicker(s.Images,
		s.Locations, s.Authorizer, s.Logger)

	site, _ := loc.NewSite("site-a", "thegroup")
	cam_a, _ := loc.NewCamera("cam-a", site, "")
	cam_b, _ := loc.NewCamera("cam-b", site, "")
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveCamera(ctx, cam_a)
	s.Locations.SaveCamera(ctx, cam_b)

	cameras := []*loc.Camera{cam_a, cam_b}

	return ctx, &s, picker, image, site, cameras
}

func TestInitialSiteAndCamera(t *testing.T) {
	ctx, s, picker, image, _, cameras := InitializeLocationPickerTests(t)
	s.Images.AssignCamera(ctx, cameras[0], image)

	state, _ := picker.Init(ctx, image)

	if state.Site == nil {
		t.Fatalf("expected state with initial site, but got none")
	}
	if state.Site.Id != cameras[0].Site.Id {
		t.Fatalf("expected state with initial site %v but got %v", cameras[0].Site, state.Site)
	}
	if state.Camera.Id != cameras[0].Id {
		t.Fatalf("expected state with initial camera %v but got %v", cameras[0], state.Camera)
	}
}

func TestAvailableSitesShouldShowAllSitesWhenOriginalIsEmpty(t *testing.T) {
	ctx, _, picker, image, _, _ := InitializeLocationPickerTests(t)

	state, _ := picker.Init(ctx, image)

	if len(state.AvailableSites) != 1 {
		t.Fatalf("expected to have 1 available site but got %v", len(state.AvailableSites))
	}
}

func TestAvailableCamerasShouldBeEmptyWhenSiteIsNotSet(t *testing.T) {
	ctx, _, picker, image, _, _ := InitializeLocationPickerTests(t)

	state, _ := picker.Init(ctx, image)

	if len(state.AvailableCameras) != 0 {
		t.Fatalf("expected to have no available cameras but got %v", len(state.AvailableCameras))
	}
}

func TestAvailableCamerasShouldBePopulatedWhenSiteIsSelected(t *testing.T) {
	ctx, _, picker, image, site, _ := InitializeLocationPickerTests(t)

	state, err := picker.SelectSite(ctx, site, image)
	AssertNoError(t, err)

	if state.Site == nil {
		t.Fatal("expected to have a selected site but got none")
	} else {
		if len(state.AvailableCameras) != 2 {
			t.Fatalf("expected to have a list of 2 cameras, but got list with length %v",
				len(state.AvailableCameras))
		}
	}
}

func TestFirstCameraShouldBePreSelectedOnSiteSelection(t *testing.T) {
	ctx, _, picker, image, site, _ := InitializeLocationPickerTests(t)

	state, err := picker.SelectSite(ctx, site, image)
	AssertNoError(t, err)

	if state.Camera == nil {
		t.Fatal("expected to pre-select a camera when a site is selected but got none ")
	}
}

func TestReselectingSiteShouldNotResetCamera(t *testing.T) {
	ctx, s, picker, image, site, cameras := InitializeLocationPickerTests(t)

	s.Images.AssignCamera(ctx, cameras[1], image)
	state, _ := picker.SelectSite(ctx, site, image)

	if state.Camera.Id != cameras[1].Id {
		t.Fatal("expected pre-selected camera to be kept as is when re-selecting original site")
	}
}

func TestCurrentCameraShouldApplyUponSelection(t *testing.T) {
	ctx, _, picker, image, _, cameras := InitializeLocationPickerTests(t)
	camera := cameras[1]
	state, _ := picker.SelectCamera(ctx, camera, image)
	if state.Camera == nil {
		t.Fatal("expected to apply a camera but did not")
	} else {
		if state.Camera.Id != camera.Id {
			t.Fatalf("expected to apply camera with id %v, but got %v", camera.Id, state.Camera.Id)
		}

	}

}
