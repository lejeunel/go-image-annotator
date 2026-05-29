package annotation

import (
	"errors"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestInternalErrOnAddLabelShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddImageLabel(a.NewAnnotationId(), image.Id, collection.Id, label.Id)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnFindImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, _ := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestAddAndRetrieveImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.AddImageLabel(a.NewAnnotationId(), image.Id, collection.Id, label.Id)
	labels, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	if err != nil {
		t.Fatalf("expected no error on find labels, got %v", err)
	}
	if len(labels) != 1 {
		t.Fatalf("expected to retrieve 1 label, got %v", len(labels))
	}
}
