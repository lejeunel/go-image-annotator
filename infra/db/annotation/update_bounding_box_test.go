package annotation

import (
	"errors"
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestInternalErrOnUpdateBoundingBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	annotationId := a.NewAnnotationId()
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, *bbox)
	repos.Annotation.Db.Close()
	err := repos.Annotation.UpdateBoundingBox(annotationId,
		a.BoundingBoxUpdatables{LabelId: label.Id, Xc: 1, Yc: 1, Width: 1, Height: 1})

	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestUpdateBoundingBoxWithInvalidValuesShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	annotationId := a.NewAnnotationId()
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, *bbox)

	err := repos.Annotation.UpdateBoundingBox(annotationId,
		a.BoundingBoxUpdatables{LabelId: label.Id, Xc: 1, Yc: 1, Width: -10, Height: 1})
	if !errors.Is(err, e.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestUpdateBoundingBox(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	annotationId := a.NewAnnotationId()
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, *bbox)

	newWidth := float32(2)
	err := repos.Annotation.UpdateBoundingBox(annotationId,
		a.BoundingBoxUpdatables{LabelId: label.Id, Xc: 1, Yc: 1, Width: newWidth, Height: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	r, _ := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	if r[0].Width != newWidth {
		t.Fatalf("expected to update width to %v, got %v", newWidth, r[0].Width)
	}
}
