package annotation

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnRemoveAnnotationShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	annotationId := a.NewAnnotationId()
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label), nil, nil)
	db.Close()
	err := repos.Annotation.RemoveAnnotation(annotationId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRemoveAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	annotationId := a.NewAnnotationId()
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label), nil, nil)
	err := repos.Annotation.RemoveAnnotation(annotationId)
	assert.NoError(t, err)
}

func TestInternalErrOnRemoveImageLabelShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	db.Close()
	err := repos.Annotation.RemoveImageLabel(image.Id, collection.Id, label.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRemoveImageLabel(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label), nil, nil)
	err := repos.Annotation.RemoveImageLabel(image.Id, collection.Id, label.Id)
	assert.NoError(t, err)
	labels, _ := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.Equal(t, 0, len(labels))
}
