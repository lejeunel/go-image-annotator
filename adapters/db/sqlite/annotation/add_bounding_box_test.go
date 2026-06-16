package annotation

import (
	"testing"
	"time"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnAddBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.Db.Close()
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, nil, nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnFindBBoxShouldFail(t *testing.T) {
	repos := NewAnnotationTestRepos()
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, nil, nil)
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
	user := u.NewUser("user@example.com")
	repos.User.Create(user)
	now := time.Now()
	err := repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, &user.Id, &now)
	assert.NoError(t, err)
	r, err := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(r))
	assert.Equal(t, labelName, r[0].Label.Name)
	assert.Equal(t, bbox.Angle, r[0].Angle)
	assert.NotNil(t, r[0].Author)
	assert.Equal(t, user.Id, *r[0].Author)
	assert.NotNil(t, r[0].Time)
}
