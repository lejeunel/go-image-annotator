package tests

import (
	"context"
	a "datahub/app"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	loc "datahub/domain/locations"
	"testing"
	"time"
)

type UpdateImageTestEnv struct {
	ctx          context.Context
	app          *a.App
	image        *im.Image
	collection   *clc.Collection
	camera       *loc.Camera
	secondCamera *loc.Camera
	otherCamera  *loc.Camera
}

func NewUpdateImageTestEnv(t *testing.T) *UpdateImageTestEnv {
	t.Helper()
	s, _, ctx := a.NewTestApp(t, true)

	image, _ := im.New(testJPGImage)
	collection, _ := clc.New("my-collection")
	s.Collections.Create(ctx, collection)
	s.Images.Save(ctx, image, collection)

	site, _ := loc.NewSite("my-site")
	camera, _ := loc.NewCamera("my-camera", site)
	secondCamera, _ := loc.NewCamera("my-second-camera", site)
	s.Locations.SaveSite(ctx, site)
	s.Locations.SaveCamera(ctx, camera)
	s.Locations.SaveCamera(ctx, secondCamera)

	otherSite, _ := loc.NewSite("my-other-site")
	otherCamera, _ := loc.NewCamera("my-other-camera", otherSite)
	s.Locations.SaveSite(ctx, otherSite)
	s.Locations.SaveCamera(ctx, otherCamera)

	return &UpdateImageTestEnv{ctx: ctx, app: &s, image: image, collection: collection,
		camera:       camera,
		secondCamera: secondCamera, otherCamera: otherCamera}

}

type UpdateImageTestCase struct {
	testName   string
	capturedAt string
	site       string
	camera     string
	type_      string
	wantErr    bool
}

func assertImageUpdated(t *testing.T, img *im.BaseImage, tc UpdateImageTestCase, layout string) {

	t.Helper()

	wantTime, _ := time.Parse(layout, tc.capturedAt)

	if !img.CapturedAt.Equal(wantTime) {
		t.Fatalf("captured_at: want %v, got %v", wantTime, img.CapturedAt)
	}

	if img.GetSiteName() != tc.site {
		t.Fatalf("site: want %v, got %v", tc.site, img.GetSiteName())
	}

	if img.GetCameraName() != tc.camera {
		t.Fatalf("camera: want %v, got %v", tc.camera, img.GetCameraName())
	}

	if img.Type != tc.type_ {
		t.Fatalf("type: want %v, got %v", tc.type_, img.Type)
	}
}

func TestUpdateImage(t *testing.T) {

	timeLayout := "2006-01-02T15:04:05.000Z"
	tests := []UpdateImageTestCase{
		{"non-existing site", "non-existing-site", "my-camera",
			"2006-01-02T15:04:05.000Z", "rgb", true},
		{"non-existing camera", "my-site", "non-existing-camera",
			"2006-01-02T15:04:05.000Z", "rgb", true},
		{"invalid timestamp", "my-site", "my-camera",
			"2006-01-02T15:04:00Z", "rgb", true},
		{"invalid type", "my-site", "my-camera",
			"2006-01-02T15:04:00Z", "asdf", true},
		{"changing type", "my-site", "my-camera",
			"2006-01-02T15:04:05.000Z", "thermal", false},
		{"changing camera", "my-other-site", "my-other-camera",
			"2006-01-02T15:04:05.000Z", "thermal", false},
	}

	for _, tc := range tests {

		env := NewUpdateImageTestEnv(t)

		updated, err := env.app.Images.Update(
			env.ctx, env.image.Id,
			im.ImageUpdatables{Site: tc.site, Camera: tc.camera, CapturedAt: tc.capturedAt, Type_: tc.type_})
		if tc.wantErr {
			if err == nil {
				t.Fatalf("%v: expected error but got nil", tc.testName)
			}
			return
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertImageUpdated(t, updated, tc, timeLayout)

	}

}
