package annotation

import (
	"testing"
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnUpdateBoundingBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	annotationId := a.NewAnnotationId()
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, nil, nil)
	repos.Annotation.Db.Close()
	err := repos.Annotation.UpdateBoundingBox(annotationId,
		a.BoundingBoxUpdatables{LabelId: label.Id, Xc: 1, Yc: 1, Width: 1, Height: 1}, nil, nil)

	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdateBoundingBoxWithInvalidValuesShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	annotationId := a.NewAnnotationId()
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, nil, nil)

	err := repos.Annotation.UpdateBoundingBox(annotationId,
		a.BoundingBoxUpdatables{LabelId: label.Id, Xc: 1, Yc: 1, Width: -10, Height: 1}, nil, nil)
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestUpdateBoundingBox(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	annotationId := a.NewAnnotationId()
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, nil, nil)
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "a-new-label")
	repos.Label.Create(newLabel)
	user := u.NewUser("user@example.com")
	repos.User.Create(user)

	newBox := a.BoundingBoxUpdatables{LabelId: newLabel.Id, Xc: 2, Yc: 3, Width: 4, Height: 10,
		Angle: -1}

	now := time.Now()
	err := repos.Annotation.UpdateBoundingBox(annotationId, newBox, &user.Id, &now)
	assert.NoError(t, err)
	r, _ := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.Equal(t, r[0].Width, newBox.Width)
	assert.Equal(t, r[0].Height, newBox.Height)
	assert.Equal(t, r[0].Xc, newBox.Xc)
	assert.Equal(t, r[0].Yc, newBox.Yc)
	assert.Equal(t, r[0].Angle, newBox.Angle)
	assert.Equal(t, r[0].Label.Id, newLabel.Id)
	assert.NotNil(t, r[0].Author)
	assert.Equal(t, user.Id, *r[0].Author)
	assert.NotNil(t, r[0].Time)
}
