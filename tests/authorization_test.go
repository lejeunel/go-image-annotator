package tests

import (
	"context"
	m "go-image-annotator/models"
	"testing"
)

func TestContributingImagesRequiresPermission(t *testing.T) {

	tests := map[string]struct {
		roles     string
		wantError bool
	}{
		"viewer role should fail":       {roles: "viewer", wantError: true},
		"admin role should succed":      {roles: "admin", wantError: false},
		"im-contrib role should succed": {roles: "im-contrib", wantError: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s, ctx := NewTestApp(t, 2)
			ctx = context.WithValue(ctx, "user_roles", tc.roles)
			image := &m.Image{Data: testImage}

			image, err := s.Images.Save(ctx, image)

			if tc.wantError {
				AssertError(t, err)

			} else {
				AssertNoError(t, err)
			}
		})
	}
}

func TestDeletingImagesRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "im-contrib")
	image := &m.Image{Data: testImage}

	image, err := s.Images.Save(ctx, image)
	err = s.Images.Delete(ctx, image)
	AssertError(t, err)
}

func TestCreatingAnnotationRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "viewer")

	label := &m.Label{Name: "mylabel"}
	label, err := s.Annotations.Create(ctx, label)
	AssertError(t, err)
}

func TestApplyingAnnotationRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "admin")

	image := &m.Image{Data: testImage}
	image, err := s.Images.Save(ctx, image)
	AssertNoError(t, err)

	label := &m.Label{Name: "mylabel"}
	label, err = s.Annotations.Create(ctx, label)
	AssertNoError(t, err)

	polyg := &m.Polygon{Label: label}

	ctx = context.WithValue(ctx, "user_roles", "im-contrib")
	image, err = s.Annotations.ApplyPolygonToImage(ctx, polyg, image)
	AssertError(t, err)

}
