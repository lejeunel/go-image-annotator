package annotation

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnUpdateLabelShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	annotationId := a.NewAnnotationId()
	label := lbl.NewLabel(lbl.NewLabelId(), "new-label")
	bbox := a.NewBoundingBox(annotationId, 1, 1, 1, 1, label)
	db.Close()
	err := repos.Annotation.UpdateLabelOfAnnotation(bbox.Id, label.Id, nil, nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdateLabelOfAnnotation(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	bbox := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, label)
	repos.Annotation.AddBoundingBox(image.Id, collection.Id, bbox, nil, nil)
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repos.Label.Create(newLabel)
	user := u.NewUser("user@example.com")
	repos.User.Create(user)
	now := time.Now()
	repos.Annotation.UpdateLabelOfAnnotation(bbox.Id, newLabel.Id, &user.Id, &now)
	r, _ := repos.Annotation.FindBoundingBoxes(image.Id, collection.Id)
	assert.Equal(t, newLabel.Id, r[0].Label.Id)
	assert.NotNil(t, r[0].Time)
	assert.NotNil(t, r[0].Author)
	assert.Equal(t, user.Id, *r[0].Author)
}
