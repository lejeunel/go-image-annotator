package tests

import (
	"context"

	m "go-image-annotator/models"
	"testing"
)

func TestSavingImagesRequiresPermission(t *testing.T) {

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
			collection := &m.Collection{Name: "myimageset"}
			s.Collections.Create(ctx, collection)
			err := s.Images.Save(ctx, image, collection)

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
	collection := &m.Collection{Name: "myimageset"}
	s.Collections.Create(ctx, collection)

	err := s.Images.Save(ctx, image, collection)
	err = s.Images.Delete(ctx, image, collection)
	AssertError(t, err)
}

func TestCreatingLabelRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)

	label := &m.Label{Name: "mylabel"}
	ctx = context.WithValue(ctx, "user_roles", "annotation-contrib")
	err := s.Annotations.CreateLabel(ctx, label)
	AssertNoError(t, err)

	ctx = context.WithValue(ctx, "user_roles", "viewer")
	err = s.Annotations.CreateLabel(ctx, label)
	AssertError(t, err)

}
func TestDeletingLabelRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)

	label := &m.Label{Name: "mylabel"}
	ctx = context.WithValue(ctx, "user_roles", "annotation-contrib")
	err := s.Annotations.CreateLabel(ctx, label)

	ctx = context.WithValue(ctx, "user_roles", "viewer")
	err = s.Annotations.DeleteLabel(ctx, label)
	AssertError(t, err)

	ctx = context.WithValue(ctx, "user_roles", "annotation-contrib")
	err = s.Annotations.DeleteLabel(ctx, label)
	AssertNoError(t, err)

}

func TestApplyingBBoxRequiresPermission(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "admin")

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: "mylabel"}
	s.Annotations.CreateLabel(ctx, label)
	collection := &m.Collection{Name: "myimageset"}
	s.Collections.Create(ctx, collection)
	err := s.Images.Save(ctx, image, collection)

	ctx = context.WithValue(ctx, "user_roles", "viewer")
	bbox := &m.BoundingBox{Xc: 10, Yc: 10, Height: 11, Width: 15}
	bbox.Annotate(label)
	err = s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collection)
	AssertError(t, err)

	ctx = context.WithValue(ctx, "user_roles", "annotation-contrib")
	err = s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collection)
	AssertNoError(t, err)

}

func TestDeletingAnnotationOnImageDoneByAnotherUserShouldFail(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "im-contrib,annotation-contrib")
	ctx = context.WithValue(ctx, "user_email", "bob@mail.com")

	image := &m.Image{Data: testImage}
	label := &m.Label{Name: "mylabel"}
	collection := &m.Collection{Name: "myimageset"}
	s.Collections.Create(ctx, collection)
	s.Annotations.CreateLabel(ctx, label)
	s.Images.Save(ctx, image, collection)

	s.Annotations.ApplyLabelToImage(ctx, label, image, collection)

	ctx = context.WithValue(ctx, "user_email", "not-bob@mail.com")
	err := s.Annotations.RemoveAnnotationFromImage(ctx, image.Annotations[0], image, collection)
	AssertError(t, err)
	if len(image.Annotations) < 1 {
		t.Fatal("expected that label is not deleted, but it is.")
	}

}

func TestDeletingBBoxDoneByAnotherUserShouldFail(t *testing.T) {

	s, ctx := NewTestApp(t, 2)
	ctx = context.WithValue(ctx, "user_roles", "im-contrib,annotation-contrib")
	ctx = context.WithValue(ctx, "user_email", "bob@mail.com")

	label := &m.Label{Name: "mylabel"}
	image := &m.Image{Data: testImage}
	collection := &m.Collection{Name: "mycollection"}
	s.Annotations.CreateLabel(ctx, label)
	s.Collections.Create(ctx, collection)

	s.Images.Save(ctx, image, collection)

	bbox := &m.BoundingBox{Xc: 10, Yc: 10, Height: 11, Width: 15}
	err := s.Annotations.ApplyBoundingBoxToImage(ctx, bbox, image, collection)
	ctx = context.WithValue(ctx, "user_email", "not-bob@mail.com")
	err = s.Annotations.RemoveAnnotationFromImage(ctx, &image.BoundingBoxes[0].Annotation, image, collection)
	AssertError(t, err)

}
