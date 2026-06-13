package annotation

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnUpdateLabelShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	annotationId := a.NewAnnotationId()
	label := lbl.NewLabel(lbl.NewLabelId(), "new-label")
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.Db.Close()
	err := repos.Annotation.UpdateLabelOfAnnotation(bbox.Id, label.Id)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdateLabelOfAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox)
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repos.Label.Create(newLabel)
	repos.Annotation.UpdateLabelOfAnnotation(bbox.Id, newLabel.Id)
	updatedBoxes, _ := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.True(t, updatedBoxes[0].Label.Id == newLabel.Id)
}
