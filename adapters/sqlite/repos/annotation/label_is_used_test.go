package annotation

import (
	"errors"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestInternalErrOnLabelIsUsedShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	imLabel := a.NewImageLabel(label)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel)
	repos.Annotation.Db.Close()
	_, err := repos.Label.IsUsed(label.Name)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestLabelIsUsedbyAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	imLabel := a.NewImageLabel(label)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel)
	isUsed, err := repos.Label.IsUsed(label.Name)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if !(*isUsed) {
		t.Fatal("expected label to be used")
	}
}
