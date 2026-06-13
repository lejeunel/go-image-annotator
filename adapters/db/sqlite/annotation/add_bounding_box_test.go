package annotation

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnAddBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddBoundingBox(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	image, collection, label := CreateAnnotableImage(repos, "a-collection", labelName, nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label,
		a.WithAngle(float32(15.)))
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	assert.NoError(t, err)
	r, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(r))
	assert.Equal(t, labelName, r[0].Label.Name)
	assert.Equal(t, bbox.Angle, r[0].Angle)
}

func TestAddImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, _ := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	newLabelName := "new-label"
	newLabel := l.NewLabel(l.NewLabelId(), newLabelName)
	imLabel := a.NewImageLabel(newLabel)
	repos.Label.Create(newLabel)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel)
	imageLabels, _ := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.Equal(t, 1, len(imageLabels))
	assert.Equal(t, newLabelName, imageLabels[0].Label.Name)

}
