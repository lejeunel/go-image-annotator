package annotation

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnAddLabelShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label))
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, _ := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddAndRetrieveImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	repos.Annotation.AddImageLabel(image.Id, collection.Id, a.NewImageLabel(label))
	labels, err := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(labels))
}
