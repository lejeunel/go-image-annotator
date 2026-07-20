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

func TestErrOnUpdateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repos := NewAnnotationTestRepos(db)
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	polygon := a.NewPolygon(a.NewAnnotationId(), TestingPolygonPoints, label)
	repos.Annotation.AddPolygon(image.Id, collection.Id, polygon, nil, nil)
	db.Close()
	err := repos.Annotation.UpdatePolygon(polygon.Id,
		a.PolygonUpdatables{LabelId: label.Id, Points: a.Points{Coordinates: [][2]float32{{0, 0}, {3, 3}}}}, nil, nil)

	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdatePolygon(t *testing.T) {
	repos := NewAnnotationTestRepos(s.NewInMemory())
	image, collection, label := CreateAnnotableImage(repos, "a-collection", "a-label", nil)
	polygon := a.NewPolygon(a.NewAnnotationId(), TestingPolygonPoints, label)
	repos.Annotation.AddPolygon(image.Id, collection.Id, polygon, nil, nil)
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "a-new-label")
	repos.Label.Create(newLabel)
	user := u.NewUser("user@example.com")
	repos.User.Create(user)

	newPolygon := a.PolygonUpdatables{LabelId: newLabel.Id, Points: a.Points{Coordinates: [][2]float32{{0, 0}, {5, 5}}}}
	now := time.Now()
	err := repos.Annotation.UpdatePolygon(polygon.Id, newPolygon, &user.Id, &now)
	assert.NoError(t, err)
	r, _ := repos.Annotation.FindPolygons(image.Id, collection.Id)
	assert.Equal(t, r[0].Label.Id, newLabel.Id)
	assert.Equal(t, r[0].Points.Coordinates, newPolygon.Points.Coordinates)
	assert.NotNil(t, r[0].Author)
	assert.Equal(t, user.Id, *r[0].Author)
	assert.NotNil(t, r[0].Time)
}
