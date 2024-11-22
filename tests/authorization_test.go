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

func TestCreatingLabelRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)

	label := &m.Label{Name: "mylabel"}
	ctx = context.WithValue(ctx, "user_roles", "annotation-contrib")
	label, err := s.Annotations.Create(ctx, label)
	AssertNoError(t, err)

	ctx = context.WithValue(ctx, "user_roles", "viewer")
	label, err = s.Annotations.Create(ctx, label)
	AssertError(t, err)

}

func TestApplyingPolygonRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "admin")

	image := &m.Image{Data: testImage}
	image, err := s.Images.Save(ctx, image)
	label := &m.Label{Name: "mylabel"}
	label, err = s.Annotations.Create(ctx, label)

	polyg := &m.Polygon{Label: label}
	ctx = context.WithValue(ctx, "user_roles", "viewer")
	image, err = s.Annotations.ApplyPolygonToImage(ctx, polyg, image)
	AssertError(t, err)

	ctx = context.WithValue(ctx, "user_roles", "annotation-contrib")
	image, err = s.Annotations.ApplyPolygonToImage(ctx, polyg, image)
	AssertNoError(t, err)

}

func TestDeletingAnnotationOnImageDoneByAnotherUserShouldFail(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "im-contrib,annotation-contrib")
	ctx = context.WithValue(ctx, "user_email", "bob@mail.com")

	label, _ := s.Annotations.Create(ctx, &m.Label{Name: "mylabel"})
	image, _ := s.Images.Save(ctx, &m.Image{Data: testImage})
	image, _ = s.Annotations.ApplyLabelToImage(ctx, label, image)

	ctx = context.WithValue(ctx, "user_email", "not-bob@mail.com")
	image, err := s.Annotations.RemoveAnnotationFromImage(ctx, image.Annotations[0], image)
	AssertError(t, err)
	if len(image.Annotations) < 1 {
		t.Fatal("expected that label is not deleted, but it is.")
	}

	ctx = context.WithValue(ctx, "user_roles", "admin")
	image, err = s.Annotations.RemoveAnnotationFromImage(ctx, image.Annotations[0], image)
	AssertNoError(t, err)

}

func TestDeletingPolygonDoneByAnotherUserShouldFail(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "im-contrib,annotation-contrib")
	ctx = context.WithValue(ctx, "user_email", "bob@mail.com")

	label, _ := s.Annotations.Create(ctx, &m.Label{Name: "mylabel"})
	image, _ := s.Images.Save(ctx, &m.Image{Data: testImage})

	polygon, err := m.NewBoundingBox(10, 10, 30, 30)
	polygon.Label = label
	image, err = s.Annotations.ApplyPolygonToImage(ctx, polygon, image)

	ctx = context.WithValue(ctx, "user_email", "not-bob@mail.com")
	image, err = s.Annotations.DeletePolygonFromImage(ctx, image.Polygons[0], image)
	AssertError(t, err)
	if len(image.Polygons) < 1 {
		t.Fatal("expected that label is not deleted, but it is.")
	}
	ctx = context.WithValue(ctx, "user_roles", "admin")
	image, err = s.Annotations.DeletePolygonFromImage(ctx, image.Polygons[0], image)
	AssertNoError(t, err)

}
