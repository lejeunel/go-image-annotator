package annotation

import (
	"errors"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestInternalErrOnRemoveAnnotationShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	annotationId := a.NewAnnotationId()
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label))
	repos.Annotation.Db.Close()
	err := repos.Annotation.RemoveAnnotation(annotationId)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestRemoveAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	annotationId := a.NewAnnotationId()
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label))
	err := repos.Annotation.RemoveAnnotation(annotationId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestInternalErrOnRemoveImageLabelShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.Db.Close()
	err := repos.Annotation.RemoveImageLabel(image.Id, collection.Id, label.Id)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestRemoveImageLabel(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label))
	err := repos.Annotation.RemoveImageLabel(image.Id, collection.Id, label.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	labels, _ := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	if len(labels) != 0 {
		t.Fatalf("expected zero image labels, got %v", len(labels))
	}
}
