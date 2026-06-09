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
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label")
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	repos.Annotation.Db.Close()
	_, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddBoundingBox(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	image, collection, label := CreateAnnotableImage(repos, "a-collection", labelName)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	assert.NoError(t, err)
	boxes, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(boxes))
	assert.Equal(t, labelName, boxes[0].Label.Name)
}

func TestRetrieveImageWithBoxesAndImageLabels(t *testing.T) {
	repos := NewAnnotationTestRepos()
	labelName := "a-label"
	image, collection, label := CreateAnnotableImage(repos, "a-collection", labelName)

	newLabelName := "new-label"
	newLabel := l.NewLabel(l.NewLabelId(), newLabelName)
	imLabel := a.NewImageLabel(newLabel)
	repos.Label.Create(newLabel)
	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel)

	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, box)

	boxes, _ := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.Equal(t, 1, len(boxes))
	imageLabels, _ := repos.Annotation.FindImageLabels(image.Id, collection.Id)
	assert.Equal(t, 1, len(imageLabels))
	assert.Equal(t, newLabelName, imageLabels[0].Label.Name)

}

// func TestRetrieveGroupOfAnnotation(t *testing.T) {
// 	repos := NewAnnotationTestRepos()
// 	labelName := "a-label"
// 	image, collection, _ := CreateAnnotableImage(repos, "a-collection", labelName)

// 	newLabelName := "new-label"
// 	newLabel := l.NewLabel(l.NewLabelId(), newLabelName)
// 	imLabel := a.NewImageLabel(newLabel)
// 	repos.Label.Create(newLabel)
// 	repos.Annotation.AddImageLabel(image.Id, collection.Id, imLabel)

// 	group, err := repos.Annotation.GroupOfAnnotation(imLabel.Id)
// 	assert.Nil(t, err)
// 	assert.Equal(t, collection.Group, *group)
// }
