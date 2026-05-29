package annotation

import (
	"errors"
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestInternalErrOnAddBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, *bbox)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnFindBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, *bbox)
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestAddBoundingBox(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	image, collection, label := CreateAnnotableImage(repos, "a-collection", labelName)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, *bbox)
	if err != nil {
		t.Fatalf("expected no error on adding bbox, got %v", err)
	}
	boxes, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	if err != nil {
		t.Fatalf("expected no error on fetching bboxes, got %v", err)
	}
	if len(boxes) != 1 {
		t.Fatalf("expected to retrieve 1 bounding box, got %v", len(boxes))
	}
	if boxes[0].Label.Name != labelName {
		t.Fatalf("expected to retrieve bbox with label name %v, got %v", labelName, boxes[0].Label.Name)
	}
}

func TestRetrieveImageWithBoxesAndImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	image, collection, label := CreateAnnotableImage(repos, "a-collection", labelName)

	newLabelName := "new-label"
	newLabel := l.NewLabel(l.NewLabelId(), newLabelName)
	repos.Label.Create(*newLabel)
	repos.Annotation.AddImageLabel(a.NewAnnotationId(), image.Id, collection.Id, newLabel.Id)

	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, *box)

	boxes, _ := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	if len(boxes) != 1 {
		t.Fatalf("expected to retrieve 1 bounding box, got %v", len(boxes))
	}
	imageLabels, _ := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	if len(imageLabels) != 1 {
		t.Fatalf("expected to retrieve 1 image label, got %v", len(imageLabels))
	}
	if imageLabels[0].Label.Name != newLabelName {
		t.Fatalf("expected to retrieve image label with name %v, got %v",
			newLabelName, imageLabels[0].Label.Name)
	}

}
